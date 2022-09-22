package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shikingram/adctl/pkg/chart/loader"
)

func Upgrade(name string) error {
	ctx := context.Background()
	return UpgradeWithContext(ctx, name)
}

var upgradeRegex = regexp.MustCompile(`^\d+-app-.*$`)

func UpgradeWithContext(ctx context.Context, name string) error {
	files, err := loader.LoadDir(filepath.Join("instance", name))
	if err != nil {
		return err
	}

	// TODO: validate

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
				return Start(fi.Name, name)
			}
		}
	}

	return nil

}
