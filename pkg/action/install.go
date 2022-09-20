package action

import (
	"auto-compose/pkg/chart"
	"auto-compose/pkg/chartutil"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const defaultDirectoryPermission = 0755

type Install struct {
	cfg *Configuration

	ReleaseName    string
	GenerateName   bool
	NameTemplate   string
	DryRun         bool
	Namespace      string
	UseReleaseName bool
}

func NewInstall(cfg *Configuration) *Install {
	return &Install{cfg: cfg}
}

func (i *Install) Name(args []string) (string, error) {
	flagsNotSet := func() error {
		if i.GenerateName {
			return errors.New("cannot set --generate-name and also specify a name")
		}
		return nil
	}

	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}

	if len(args) == 1 {
		return args[0], flagsNotSet()
	}

	if i.ReleaseName != "" {
		return i.ReleaseName, nil
	}

	if !i.GenerateName {
		return "", errors.New("must either provide a name or specify --generate-name")
	}

	base := filepath.Base(args[0])
	if base == "." || base == "" {
		base = "chart"
	}

	if idx := strings.Index(base, "."); idx != -1 {
		base = base[0:idx]
	}

	return fmt.Sprintf("%s-%d", base, time.Now().Unix()), nil
}

func writeToFile(dir string, sourceName string, data string, append bool) error {
	outfileName := strings.Join([]string{dir, sourceName}, string(filepath.Separator))

	err := ensureDirectoryForFile(outfileName)
	if err != nil {
		return err
	}

	f, err := createOrOpenFile(outfileName, append)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("---\n# Source: %s\n%s\n", sourceName, data))

	if err != nil {
		return err
	}

	fmt.Printf("wrote %s\n", outfileName)
	return nil
}

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(filename)
}

// check if the directory exists to create file. creates if don't exists
func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, defaultDirectoryPermission)
}

func (i *Install) Run() {
}

func (i *Install) RunWithContext(ctx context.Context, ch *chart.Chart, vals chartutil.Values) error {
	options := chartutil.ReleaseOptions{
		Name:     i.ReleaseName,
		Revision: 1,
	}

	valuesToRender, err := chartutil.ToRenderValues(ch, options, vals)
	if err != nil {
		return err
	}
	return i.cfg.renderResources(ch, valuesToRender, i.ReleaseName)
}
