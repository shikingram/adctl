package values

import (
	"io/ioutil"

	"github.com/shikingram/auto-compose/pkg/chartutil"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

type Options struct {
	ValueFiles   []string
	StringValues []string
	Values       []string
	FileValues   []string
}

// MergeValues merges values from files specified via -f/--values and directly
// via --set  marshaling them to YAML
func (opts *Options) MergeValues() (chartutil.Values, error) {
	base := chartutil.Values{}

	// // local default values files
	// defaultYaml, err := ioutil.ReadFile("values.yaml")
	// if err != nil {
	// 	return nil, err
	// }
	// if err := yaml.Unmarshal(defaultYaml, &base); err != nil {
	// 	return nil, errors.Wrapf(err, "failed to parse default yaml")
	// }

	// User specified a values files via -f/--values
	for _, filePath := range opts.ValueFiles {
		currentMap := chartutil.Values{}

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s", filePath)
		}

		// Merge with the previous map
		base = mergeMaps(base, currentMap)
	}

	// TODO: User specified a value via --set

	return base, nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
