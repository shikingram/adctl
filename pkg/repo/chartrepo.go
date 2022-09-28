package repo

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/pkg/adctlpath"
	"github.com/shikingram/adctl/pkg/getter"
	"sigs.k8s.io/yaml"
)

// Entry represents a collection of parameters for chart repository
type Entry struct {
	Name                  string `json:"name"`
	URL                   string `json:"url"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	CertFile              string `json:"certFile"`
	KeyFile               string `json:"keyFile"`
	CAFile                string `json:"caFile"`
	InsecureSkipTLSverify bool   `json:"insecure_skip_tls_verify"`
	PassCredentialsAll    bool   `json:"pass_credentials_all"`
}

type ChartRepository struct {
	Config     *Entry
	ChartPaths []string
	Client     getter.Getter
	IndexFile  *IndexFile
	CachePath  string
}

func NewChartRepository(cfg *Entry, getters getter.Providers) (*ChartRepository, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, errors.Errorf("invalid chart URL format: %s", cfg.URL)
	}

	client, err := getters.ByScheme(u.Scheme)
	if err != nil {
		return nil, errors.Errorf("could not find protocol handler for: %s", u.Scheme)
	}

	return &ChartRepository{
		Config:    cfg,
		IndexFile: NewIndexFile(),
		Client:    client,
		CachePath: adctlpath.CachePath("repository"),
	}, nil
}

func (r *ChartRepository) DownloadIndexFile() (string, error) {
	parsedURL, err := url.Parse(r.Config.URL)
	if err != nil {
		return "", err
	}
	parsedURL.RawPath = path.Join(parsedURL.RawPath, "index.yaml")
	parsedURL.Path = path.Join(parsedURL.Path, "index.yaml")

	indexURL := parsedURL.String()
	resp, err := r.Client.Get(indexURL,
		getter.WithURL(r.Config.URL),
		getter.WithBasicAuth(r.Config.Username, r.Config.Password),
	)
	if err != nil {
		return "", err
	}

	index, err := ioutil.ReadAll(resp)
	if err != nil {
		return "", err
	}

	indexFile, err := loadIndex(index, r.Config.URL)
	if err != nil {
		return "", err
	}

	// Create the chart list file in the cache directory
	var charts strings.Builder
	for name := range indexFile.Entries {
		fmt.Fprintln(&charts, name)
	}
	chartsFile := filepath.Join(r.CachePath, adctlpath.CacheChartsFile(r.Config.Name))
	os.MkdirAll(filepath.Dir(chartsFile), 0755)
	ioutil.WriteFile(chartsFile, []byte(charts.String()), 0644)

	// Create the index file in the cache directory
	fname := filepath.Join(r.CachePath, adctlpath.CacheIndexFile(r.Config.Name))
	os.MkdirAll(filepath.Dir(fname), 0755)
	index, err = yaml.Marshal(indexFile)
	if err != nil {
		return "", err
	}
	return fname, ioutil.WriteFile(fname, index, 0644)
}
