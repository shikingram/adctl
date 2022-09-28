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

func Install(name string, ch *chart.Chart) error {
	ctx := context.Background()
	return InstallWithContext(ctx, ch, name)
}

var nameRegex = regexp.MustCompile(`^\d+-(app|job)-.*$`)

const notesFileSuffix = "NOTES.txt"

func InstallWithContext(ctx context.Context, ch *chart.Chart, name string) error {
	files, err := loader.LoadDir(filepath.Join("instance", name, ch.Metadata.Name))
	if err != nil {
		return err
	}

	// validate
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
			if nameRegex.MatchString(filepath.Base(fi.Name)) {
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
			fmt.Printf("\n %s: \n %s \n", notesFileSuffix, string(fi.Data))
		}
	}
}

func validateFiles(releaseName string, files []*loader.BufferedFile, ch *chart.Chart) error {
	for _, fi := range files {
		basename := filepath.Base(fi.Name)
		if nameRegex.MatchString(basename) {
			if !infiles(basename, ch.Templates) {
				Down(fi.Name, releaseName)
				os.RemoveAll(fi.Name)
			}
		}
	}
	return nil
}

func infiles(file string, ts []*chart.File) bool {
	for _, t := range ts {
		if strings.Contains(filepath.Base(t.Name), file) {
			return true
		}
	}
	return false
}
