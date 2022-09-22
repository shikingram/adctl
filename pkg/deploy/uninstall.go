package deploy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shikingram/auto-compose/pkg/chart/loader"
)

func UnInstall(name string) error {
	rootPath := filepath.Join("instance", name)
	files, err := loader.LoadDir(rootPath)
	if err != nil {
		return err
	}
	for _, fi := range files {
		if nameRegex.MatchString(fi.Name[strings.LastIndex(fi.Name, string(os.PathSeparator))+1:]) {
			err := Down(fi.Name, name)
			if err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(rootPath)

}
