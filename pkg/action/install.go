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
	"github.com/shikingram/adctl/pkg/cli"
	"github.com/shikingram/adctl/pkg/deploy"
	"github.com/shikingram/adctl/pkg/downloader"
	"github.com/shikingram/adctl/pkg/getter"

	"github.com/pkg/errors"
)

const defaultDirectoryPermission = 0755

type Install struct {
	cfg *Configuration

	ChartPathOptions

	ReleaseName    string
	GenerateName   bool
	NameTemplate   string
	DryRun         bool
	Force          bool
	UseReleaseName bool
}

type ChartPathOptions struct {
	Version  string // --version
	RepoURL  string // --repo
	Username string // --username
	Password string // --password
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

	return deploy.InstallWithContext(ctx, ch, i.ReleaseName)
}

func (i *Install) ValidateName(name string) bool {
	num, err := deploy.CheckReleaseDeploy(name)
	return err == nil && num > 0
}

// LocateChart returns a filename of tgz
func (c *ChartPathOptions) LocateChart(name string, settings *cli.EnvSettings) (string, error) {
	name = strings.TrimSpace(name)
	version := strings.TrimSpace(c.Version)

	if _, err := os.Stat(name); err == nil {
		abs, err := filepath.Abs(name)
		if err != nil {
			return abs, err
		}

		return abs, nil
	}

	if filepath.IsAbs(name) || strings.HasPrefix(name, ".") {
		return name, errors.Errorf("path %q not found", name)
	}

	dl := downloader.ChartDownloader{
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}
	dl.Options = append(dl.Options, getter.WithBasicAuth(c.Username, c.Password))

	if err := os.MkdirAll(settings.RepositoryCache, 0755); err != nil {
		return "", err
	}

	filename, err := dl.DownloadTo(name, version, settings.RepositoryCache)
	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return lname, nil
	}

	return "", err

}
