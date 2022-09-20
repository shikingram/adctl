package loader

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/shikingram/auto-compose/pkg/chart"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

type BufferedFile struct {
	Name string
	Data []byte
}

func loadDir(path string) ([]*BufferedFile, error) {
	var files []*BufferedFile

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		fullname := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			temps, err := loadDir(fullname)
			if err != nil {
				return files, err
			}
			files = append(files, temps...)
		} else {
			data, err := ioutil.ReadFile(fullname)
			if err != nil {
				return files, err
			}
			files = append(files, &BufferedFile{Name: fullname, Data: data})
		}
	}
	return files, nil
}

//LoadChart load chart from local path
func LoadChart(path string) (*chart.Chart, error) {
	c := new(chart.Chart)

	if path == "" {
		path = "chart"
	}
	files, err := loadDir(path)
	if err != nil {
		return c, err
	}
	for _, f := range files {
		c.Raw = append(c.Raw, &chart.File{Name: f.Name, Data: f.Data})
		if strings.Contains(f.Name, "Chart.yaml") {
			if c.Metadata == nil {
				c.Metadata = new(chart.Metadata)
			}
			if err := yaml.Unmarshal(f.Data, c.Metadata); err != nil {
				return c, errors.Wrap(err, "cannot load Chart.yaml")
			}
		}
	}

	for _, f := range files {
		switch {
		case strings.Contains(f.Name, "values.yaml"):
			c.Values = make(map[string]interface{})
			if err := yaml.Unmarshal(f.Data, &c.Values); err != nil {
				return c, errors.Wrap(err, "cannot load values.yaml")
			}
		case strings.Contains(f.Name, ".gtpl"):
			c.Templates = append(c.Templates, &chart.File{Name: f.Name, Data: f.Data})
		default:
			c.Files = append(c.Files, &chart.File{Name: f.Name, Data: f.Data})
		}
	}

	if c.Metadata == nil {
		return c, errors.New("Chart.yaml file is missing")
	}

	if err := c.Validate(); err != nil {
		return c, err
	}
	return c, nil
}
