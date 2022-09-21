package chartutil

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/shikingram/auto-compose/pkg/chart"

	"sigs.k8s.io/yaml"
)

// Values represents a collection of chart values.
type Values map[string]interface{}

// YAML encodes the Values into a YAML string.
func (v Values) YAML() (string, error) {
	b, err := yaml.Marshal(v)
	return string(b), err
}

// AsMap is a utility function for converting Values to a map[string]interface{}.
//
// It protects against nil map panics.
func (v Values) AsMap() map[string]interface{} {
	if v == nil || len(v) == 0 {
		return map[string]interface{}{}
	}
	return v
}

// Encode writes serialized Values information to the given io.Writer.
func (v Values) Encode(w io.Writer) error {
	out, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

func (v Values) PathValue(path string) (interface{}, error) {
	if path == "" {
		return nil, errors.New("YAML path cannot be empty")
	}
	return v.pathValue(parsePath(path))
}

func (v Values) Table(name string) (Values, error) {
	table := v
	var err error

	for _, n := range parsePath(name) {
		if table, err = tableLookup(table, n); err != nil {
			break
		}
	}
	return table, err
}

func tableLookup(v Values, simple string) (Values, error) {
	v2, ok := v[simple]
	if !ok {
		return v, ErrNoTable{simple}
	}
	if vv, ok := v2.(map[string]interface{}); ok {
		return vv, nil
	}

	// This catches a case where a value is of type Values, but doesn't (for some
	// reason) match the map[string]interface{}. This has been observed in the
	// wild, and might be a result of a nil map of type Values.
	if vv, ok := v2.(Values); ok {
		return vv, nil
	}

	return Values{}, ErrNoTable{simple}
}

func (v Values) pathValue(path []string) (interface{}, error) {
	if len(path) == 1 {
		// if exists must be root key not table
		if _, ok := v[path[0]]; ok && !istable(v[path[0]]) {
			return v[path[0]], nil
		}
		return nil, ErrNoValue{path[0]}
	}

	key, path := path[len(path)-1], path[:len(path)-1]
	// get our table for table path
	t, err := v.Table(joinPath(path...))
	if err != nil {
		return nil, ErrNoValue{key}
	}
	// check table for key and ensure value is not a table
	if k, ok := t[key]; ok && !istable(k) {
		return k, nil
	}
	return nil, ErrNoValue{key}
}

func parsePath(key string) []string { return strings.Split(key, ".") }

func joinPath(path ...string) string { return strings.Join(path, ".") }

// ReadValues will parse YAML byte data into a Values.
func ReadValues(data []byte) (vals Values, err error) {
	err = yaml.Unmarshal(data, &vals)
	if len(vals) == 0 {
		vals = Values{}
	}
	return vals, err
}

// ReadValuesFile will parse a YAML file into a map of values.
func ReadValuesFile(filename string) (Values, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return ReadValues(data)
}

type ReleaseOptions struct {
	Name     string
	Revision int
}

func ToRenderValues(chrt *chart.Chart, options ReleaseOptions, chrtVals map[string]interface{}) (Values, error) {

	top := map[string]interface{}{
		"Chart":        chrt.Metadata,
		"Capabilities": nil,
		"Release": map[string]interface{}{
			"Name":     options.Name,
			"Revision": options.Revision,
			"Service":  "auto-compose",
		},
	}

	vals, err := CoalesceValues(chrt, chrtVals)
	if err != nil {
		return top, err
	}

	// TODO: validate
	// if err := ValidateAgainstSchema(chrt, vals); err != nil {
	// 	errFmt := "values don't meet the specifications of the schema(s) in the following chart(s):\n%s"
	// 	return top, fmt.Errorf(errFmt, err.Error())
	// }

	top["Values"] = vals
	return top, nil
}
