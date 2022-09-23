package engine

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"sigs.k8s.io/yaml"
)

func TestRender(t *testing.T) {
	t1 := template.New("gotpl")

	content, _ := ioutil.ReadFile("test.yaml")
	t1.New("test.yaml").Parse(string(content))

	m := make(map[string]interface{})

	vals, _ := ioutil.ReadFile("values.yaml")

	_ = yaml.Unmarshal(vals, &m)

	var buf strings.Builder
	_ = t1.ExecuteTemplate(&buf, "test.yaml", m)

	fmt.Println(buf.String())
}
