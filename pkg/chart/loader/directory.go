package loader

import "github.com/shikingram/adctl/pkg/chart"

type DirLoader string

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func (l DirLoader) Load() (*chart.Chart, error) {
	files, err := LoadDir(string(l))
	if err != nil {
		return nil, err
	}

	return LoadFiles(files)
}
