package repo

import (
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/shikingram/adctl/pkg/chart"
	"sigs.k8s.io/yaml"
)

var (
	ErrNoAPIVersion   = errors.New("no API version specified")
	ErrNoChartVersion = errors.New("no chart version found")
	ErrNoChartName    = errors.New("no chart name found")
	ErrEmptyIndexYaml = errors.New("empty index.yaml file")
)

type IndexFile struct {
	ServerInfo  map[string]interface{}   `json:"serverInfo,omitempty"`
	APIVersion  string                   `json:"apiVersion"`
	Generated   time.Time                `json:"generated"`
	Entries     map[string]ChartVersions `json:"entries"`
	PublicKeys  []string                 `json:"publicKeys,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
}

const APIVersionV1 = "v1"

func NewIndexFile() *IndexFile {
	return &IndexFile{
		Generated:  time.Now(),
		APIVersion: APIVersionV1,
		Entries:    map[string]ChartVersions{},
		PublicKeys: []string{},
	}
}

func (i IndexFile) SortEntries() {
	for _, versions := range i.Entries {
		sort.Sort(sort.Reverse(versions))
	}
}

func (i IndexFile) Get(name, version string) (*ChartVersion, error) {
	vs, ok := i.Entries[name]
	if !ok {
		return nil, ErrNoChartName
	}
	if len(vs) == 0 {
		return nil, ErrNoChartVersion
	}

	var constraint *semver.Constraints
	if version == "" {
		constraint, _ = semver.NewConstraint("*")
	} else {
		var err error
		constraint, err = semver.NewConstraint(version)
		if err != nil {
			return nil, err
		}
	}

	if len(version) != 0 {
		for _, ver := range vs {
			if version == ver.Version {
				return ver, nil
			}
		}
	}

	for _, ver := range vs {
		test, err := semver.NewVersion(ver.Version)
		if err != nil {
			continue
		}

		if constraint.Check(test) {
			return ver, nil
		}
	}
	return nil, errors.Errorf("no chart version found for %s-%s", name, version)
}

func LoadIndexFile(path string) (*IndexFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i, err := loadIndex(b, path)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading %s", path)
	}
	return i, nil
}

type ChartVersion struct {
	*chart.Metadata
	URLs                    []string  `json:"urls"`
	Created                 time.Time `json:"created,omitempty"`
	Removed                 bool      `json:"removed,omitempty"`
	Digest                  string    `json:"digest,omitempty"`
	ChecksumDeprecated      string    `json:"checksum,omitempty"`
	EngineDeprecated        string    `json:"engine,omitempty"`
	TillerVersionDeprecated string    `json:"tillerVersion,omitempty"`
	URLDeprecated           string    `json:"url,omitempty"`
}

type ChartVersions []*ChartVersion

func (c ChartVersions) Len() int { return len(c) }

func (c ChartVersions) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c ChartVersions) Less(a, b int) bool {
	i, err := semver.NewVersion(c[a].Version)
	if err != nil {
		return true
	}
	j, err := semver.NewVersion(c[b].Version)
	if err != nil {
		return false
	}
	return i.LessThan(j)
}

func loadIndex(data []byte, source string) (*IndexFile, error) {
	i := &IndexFile{}

	if len(data) == 0 {
		return i, ErrEmptyIndexYaml
	}

	if err := yaml.UnmarshalStrict(data, i); err != nil {
		return i, err
	}

	for name, cvs := range i.Entries {
		for idx := len(cvs) - 1; idx >= 0; idx-- {

			if !validateAnnotations(cvs[idx].Annotations) {
				continue
			}

			if cvs[idx].APIVersion == "" {
				cvs[idx].APIVersion = chart.APIVersionV1
			}

			if err := cvs[idx].Validate(); err != nil {
				log.Printf("skipping loading invalid entry for chart %q %q from %s: %s", name, cvs[idx].Version, source, err)
				cvs = append(cvs[:idx], cvs[idx+1:]...)
			}
		}
	}
	i.SortEntries()

	if i.APIVersion == "" {
		return i, ErrNoAPIVersion
	}

	return i, nil
}

func validateAnnotations(m map[string]string) bool {
	return m["category"] == "docker-compose"
}
