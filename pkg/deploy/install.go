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

func Install(name string) error {
	ctx := context.Background()
	return InstallWithContext(ctx, name)
}

var nameRegex = regexp.MustCompile(`^\d+-(app|job)-.*$`)

const notesFileSuffix = "NOTES.txt"

func InstallWithContext(ctx context.Context, name string) error {
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
			if nameRegex.MatchString(fi.Name[strings.LastIndex(fi.Name, string(os.PathSeparator))+1:]) {
				fmt.Printf("match file: %s \n", fi.Name)
				err := Start(fi.Name, name)
				if err != nil {
					return err
				}
			}
		}
	}

	printNotes(files)

	return nil

}

func printNotes(fis []*loader.BufferedFile) {
	for _, fi := range fis {
		if strings.Contains(fi.Name, notesFileSuffix) {
			fmt.Printf("\n %s ==> \n %s \n", notesFileSuffix, string(fi.Data))
		}
	}
}
