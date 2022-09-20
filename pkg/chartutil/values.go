package chartutil

import (
	"io"
	"io/ioutil"

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
