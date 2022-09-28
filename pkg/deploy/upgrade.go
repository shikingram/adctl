package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shikingram/adctl/pkg/chart"
	"github.com/shikingram/adctl/pkg/chart/loader"
)

func Upgrade(name string, ch *chart.Chart) error {
	ctx := context.Background()
	return UpgradeWithContext(ctx, ch, name, false)
}

var upgradeRegex = regexp.MustCompile(`^\d+-app-.*$`)

func UpgradeWithContext(ctx context.Context, ch *chart.Chart, name string, force bool) error {
	files, err := loader.LoadDir(filepath.Join("instance", name))
	if err != nil {
		return err
	}

	err = validateFiles(name, files, ch)
	if err != nil {
		return err
	}

	if err := CreateNetWork(name); err != nil {
		return err
	}

	for _, fi := range files {
		select {
		case <-ctx.Done():
			return errors.New("ctx canceled")
		default:
			if upgradeRegex.MatchString(fi.Name[strings.LastIndex(fi.Name, string(os.PathSeparator))+1:]) {
				fmt.Printf("match file: %s \n", fi.Name)

				var err error
				if force {
					err = Restart(fi.Name, name)
				} else {
					err = Start(fi.Name, name)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	printNotes(files)

	return nil

}
