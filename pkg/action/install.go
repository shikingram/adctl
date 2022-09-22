package action

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/shikingram/adctl/pkg/chart"
	"github.com/shikingram/adctl/pkg/chartutil"
	"github.com/shikingram/adctl/pkg/deploy"

	"github.com/pkg/errors"
)

const defaultDirectoryPermission = 0755

type Install struct {
	cfg *Configuration

	ReleaseName    string
	GenerateName   bool
	NameTemplate   string
	DryRun         bool
	UseReleaseName bool
}

func NewInstall(cfg *Configuration) *Install {
	return &Install{cfg: cfg}
}

func (i *Install) NameAndChart(args []string) (string, string, error) {
	flagsNotSet := func() error {
		if i.GenerateName {
			return errors.New("cannot set --generate-name and also specify a name")
		}
		return nil
	}

	if len(args) > 2 {
		return args[0], args[1], errors.Errorf("expected at most two arguments, unexpected arguments: %v", strings.Join(args[2:], ", "))
	}

	if len(args) == 2 {
		return args[0], args[1], flagsNotSet()
	}

	if i.ReleaseName != "" {
		return i.ReleaseName, args[0], nil
	}

	if !i.GenerateName {
		return "", args[0], errors.New("must either provide a name or specify --generate-name")
	}

	base := filepath.Base(args[0])
	if base == "." || base == "" {
		base = "chart"
	}

	if idx := strings.Index(base, "."); idx != -1 {
		base = base[0:idx]
	}

	return fmt.Sprintf("%s-%d", base, time.Now().Unix()), args[0], nil
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

	_, err = f.WriteString(data)

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

func (i *Install) Run(ch *chart.Chart, vals chartutil.Values) error {
	ctx := context.Background()
	return i.RunWithContext(ctx, ch, vals)
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
	err = i.cfg.renderResources(ch, valuesToRender, i.ReleaseName)
	if err != nil {
		return err
	}

	if i.DryRun {
		return nil
	}

	return deploy.InstallWithContext(ctx, i.ReleaseName)
}

func (i *Install) ValidateName(name string) bool {
	num, err := deploy.CheckReleaseDeploy(name)
	return err == nil && num > 0
}
