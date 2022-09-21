package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shikingram/auto-compose/pkg/chart/loader"
)

func Install(name string) error {
	ctx := context.Background()
	return InstallWithContext(ctx, name)
}

var nameRegex = regexp.MustCompile(`^\d{2}-(app|job)-[a-z]+(.yaml)$`)

func InstallWithContext(ctx context.Context, name string) error {
	files, err := loader.LoadDir("instance")
	if err != nil {
		return err
	}

	// TODO: validate

	err = CreateNetWork(name)
	if err != nil {
		return err
	}

	for _, fi := range files {
		select {
		case <-ctx.Done():
			return errors.New("ctx canceled")
		default:
			if nameRegex.MatchString(fi.Name[strings.LastIndex(fi.Name, string(os.PathSeparator))+1:]) {
				fmt.Printf("match file: %s \n", fi.Name)
				err = Start(fi.Name)
				if err != nil {
					fmt.Fprintf(os.Stderr, "start docker container failed :%s", err.Error())
					continue
				}
			}
		}
	}

	return nil

}
