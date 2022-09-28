package action

import (
	"path/filepath"

	"github.com/shikingram/adctl/pkg/chart"
	"github.com/shikingram/adctl/pkg/chartutil"
	"github.com/shikingram/adctl/pkg/engine"
	"sigs.k8s.io/yaml"
)

type Configuration struct {
}

func (cfg *Configuration) renderResources(ch *chart.Chart, values chartutil.Values, releaseName string) error {
	files, err := engine.Render(ch, values)
	if err != nil {
		return err
	}

	// add coalesced values.yaml to target directory
	content, err := yaml.Marshal(values["Values"])
	if err != nil {
		return err
	}
	files[filepath.Join(ch.ChartPath(), "values.yaml")] = string(content)

	// add Chart.yaml to target directory
	chartcontent, err := yaml.Marshal(ch.Metadata)
	if err != nil {
		return err
	}
	files[filepath.Join(ch.ChartPath(), "Chart.yaml")] = string(chartcontent)

	fileWritten := make(map[string]bool)
	for name, content := range files {
		newDir := filepath.Join("instance", releaseName)
		err = writeToFile(newDir, name, content, fileWritten[name])
		if err != nil {
			return err
		}
		fileWritten[name] = true
	}

	return nil
}
