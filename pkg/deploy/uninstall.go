package deploy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shikingram/adctl/pkg/chart/loader"
)

func UnInstall(name string, remove bool) error {
	rootpath := filepath.Join("instance", name)
	files, err := loader.LoadDir(rootpath)
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

	if remove {
		return os.RemoveAll(rootpath)
	}
	return nil

}
