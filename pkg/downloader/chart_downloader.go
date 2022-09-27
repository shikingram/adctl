package downloader

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/internal/fileutil"
	"github.com/shikingram/adctl/pkg/adctlpath"
	"github.com/shikingram/adctl/pkg/getter"
	"github.com/shikingram/adctl/pkg/repo"
)

func VerifyChart() {

}

type ChartDownloader struct {
	Out              io.Writer
	Keyring          string
	RepositoryConfig string
	RepositoryCache  string

	Getters getter.Providers
	Options []getter.Option
}

func (c *ChartDownloader) DownloadTo(ref, version, dest string) (string, error) {
	u, err := c.ResolveChartVersion(ref, version)
	if err != nil {
		return "", err
	}

	g, err := c.Getters.ByScheme(u.Scheme)
	if err != nil {
		return "", err
	}

	data, err := g.Get(u.String(), c.Options...)
	if err != nil {
		return "", err
	}

	name := filepath.Base(u.Path)
	destfile := filepath.Join(dest, name)
	if err := fileutil.AtomicWriteFile(destfile, data, 0644); err != nil {
		return destfile, err
	}

	return destfile, nil
}

var ErrNoOwnerRepo = errors.New("could not find a repo containing the given URL")

func (c *ChartDownloader) ResolveChartVersion(ref, version string) (*url.URL, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return u, err
	}

	rf, err := loadRepoConfig(c.RepositoryConfig)
	if err != nil {
		return u, err
	}

	p := strings.SplitN(u.Path, "/", 2)
	if len(p) < 2 {
		return u, errors.Errorf("non-absolute URLs should be in form of repo_name/path_to_chart, got: %s", u)
	}

	repoName := p[0]
	chartName := p[1]
	rc, err := pickChartRepositoryConfigByName(repoName, rf.Repositories)

	if err != nil {
		return u, err
	}

	c.Options = append(c.Options, getter.WithURL(rc.URL))

	r, err := repo.NewChartRepository(rc, c.Getters)
	if err != nil {
		return u, err
	}

	if r != nil && r.Config != nil {
		if r.Config.Username != "" && r.Config.Password != "" {
			c.Options = append(c.Options,
				getter.WithBasicAuth(r.Config.Username, r.Config.Password),
			)
		}
	}

	idxFile := filepath.Join(c.RepositoryCache, adctlpath.CacheIndexFile(r.Config.Name))
	i, err := repo.LoadIndexFile(idxFile)
	if err != nil {
		return u, errors.Wrap(err, "no cached repo found. (try 'adctl repo update')")
	}
	cv, err := i.Get(chartName, version)
	if err != nil {
		return u, errors.Wrapf(err, "chart %q matching %s not found in %s index. (try 'adctl repo update')", chartName, version, r.Config.Name)
	}

	if len(cv.URLs) == 0 {
		return u, errors.Errorf("chart %q has no downloadable URLs", ref)
	}

	u, err = url.Parse(cv.URLs[0])
	if err != nil {
		return u, errors.Errorf("invalid chart URL format: %s", ref)
	}

	if !u.IsAbs() {
		repoURL, err := url.Parse(rc.URL)
		if err != nil {
			return repoURL, err
		}
		q := repoURL.Query()
		repoURL.Path = strings.TrimSuffix(repoURL.Path, "/") + "/"
		u = repoURL.ResolveReference(u)
		u.RawQuery = q.Encode()
		if _, err := getter.NewHTTPGetter(getter.WithURL(rc.URL)); err != nil {
			return repoURL, err
		}
		return u, err
	}

	return u, err
}

func loadRepoConfig(file string) (*repo.File, error) {
	r, err := repo.LoadFile(file)
	if err != nil && !os.IsNotExist(errors.Cause(err)) {
		return nil, err
	}
	return r, nil
}

func pickChartRepositoryConfigByName(name string, cfgs []*repo.Entry) (*repo.Entry, error) {
	for _, rc := range cfgs {
		if rc.Name == name {
			if rc.URL == "" {
				return nil, errors.Errorf("no URL found for repository %s", name)
			}
			return rc, nil
		}
	}
	return nil, errors.Errorf("repo %s not found", name)
}
