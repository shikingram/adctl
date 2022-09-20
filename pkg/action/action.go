package action

import (
	"path/filepath"

	"github.com/shikingram/auto-compose/pkg/chart"
	"github.com/shikingram/auto-compose/pkg/chartutil"
	"github.com/shikingram/auto-compose/pkg/engine"
)

type Configuration struct {
	// specify path of application
	PrivateCfg string
}

func (cfg *Configuration) renderResources(ch *chart.Chart, values chartutil.Values, releaseName string) error {
	files, err := engine.Render(ch, values)
	if err != nil {
		return err
	}

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
