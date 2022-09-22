package engine

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/shikingram/adctl/pkg/chart"

	"github.com/shikingram/adctl/pkg/chartutil"

	"github.com/pkg/errors"
)

const warnStartDelim = "HELM_ERR_START"
const warnEndDelim = "HELM_ERR_END"
const recursionMaxNums = 1000

var warnRegex = regexp.MustCompile(warnStartDelim + `((?s).*)` + warnEndDelim)

type Engine struct {
}

// Render takes a chart, optional values, and value overrides, and attempts to
// render the Go templates using the default options.
func Render(chrt *chart.Chart, values chartutil.Values) (map[string]string, error) {
	return new(Engine).Render(chrt, values)
}

func (e Engine) Render(chrt *chart.Chart, values chartutil.Values) (map[string]string, error) {
	tmap := allTemplates(chrt, values)
	return e.render(tmap)
}

type renderable struct {
	// tpl is the current template.
	tpl string
	// vals are the values to be supplied to the template.
	vals chartutil.Values
	// namespace prefix to the templates of the current chart
	basePath string
}

func (e Engine) initFunMap(t *template.Template, referenceTpls map[string]renderable) {
	funcMap := funcMap()
	includedNames := make(map[string]int)

	// Add the 'include' function here so we can close over t.
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", errors.Wrapf(fmt.Errorf("unable to execute template"), "rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := t.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}

	// Add the 'tpl' function here
	funcMap["tpl"] = func(tpl string, vals chartutil.Values) (string, error) {
		basePath, err := vals.PathValue("Template.BasePath")
		if err != nil {
			return "", errors.Wrapf(err, "cannot retrieve Template.Basepath from values inside tpl function: %s", tpl)
		}

		templateName, err := vals.PathValue("Template.Name")
		if err != nil {
			return "", errors.Wrapf(err, "cannot retrieve Template.Name from values inside tpl function: %s", tpl)
		}

		templates := map[string]renderable{
			templateName.(string): {
				tpl:      tpl,
				vals:     vals,
				basePath: basePath.(string),
			},
		}

		result, err := e.renderWithReferences(templates, referenceTpls)
		if err != nil {
			return "", errors.Wrapf(err, "error during tpl function execution for %q", tpl)
		}
		return result[templateName.(string)], nil
	}

	t.Funcs(funcMap)
}

func (e Engine) renderWithReferences(tpls, referenceTpls map[string]renderable) (rendered map[string]string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("rendering template failed: %v", r)
		}
	}()
	t := template.New("gotpl")

	e.initFunMap(t, referenceTpls)

	keys := sortTemplates(tpls)
	for _, filename := range keys {
		r := tpls[filename]
		if _, err := t.New(filename).Parse(r.tpl); err != nil {
			return map[string]string{}, cleanupParseError(filename, err)
		}
	}

	rendered = make(map[string]string, len(keys))
	for _, filename := range keys {
		if strings.HasPrefix(path.Base(filename), "_") {
			continue
		}
		vals := tpls[filename].vals
		vals["Template"] = chartutil.Values{"Name": filename, "BasePath": tpls[filename].basePath}
		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, filename, vals); err != nil {
			return map[string]string{}, cleanupExecError(filename, err)
		}
		filename = filename[:strings.LastIndex(filename, ".gtpl")]
		rendered[filename] = strings.ReplaceAll(buf.String(), "<no value>", "")
	}

	return rendered, nil
}

func (e Engine) render(tpls map[string]renderable) (rendered map[string]string, err error) {
	return e.renderWithReferences(tpls, tpls)
}

func cleanupParseError(filename string, err error) error {
	tokens := strings.Split(err.Error(), ": ")
	if len(tokens) == 1 {
		// This might happen if a non-templating error occurs
		return fmt.Errorf("parse error in (%s): %s", filename, err)
	}
	// The first token is "template"
	// The second token is either "filename:lineno" or "filename:lineNo:columnNo"
	location := tokens[1]
	// The remaining tokens make up a stacktrace-like chain, ending with the relevant error
	errMsg := tokens[len(tokens)-1]
	return fmt.Errorf("parse error at (%s): %s", string(location), errMsg)
}

func cleanupExecError(filename string, err error) error {
	if _, isExecError := err.(template.ExecError); !isExecError {
		return err
	}

	tokens := strings.SplitN(err.Error(), ": ", 3)
	if len(tokens) != 3 {
		// This might happen if a non-templating error occurs
		return fmt.Errorf("execution error in (%s): %s", filename, err)
	}

	// The first token is "template"
	// The second token is either "filename:lineno" or "filename:lineNo:columnNo"
	location := tokens[1]

	parts := warnRegex.FindStringSubmatch(tokens[2])
	if len(parts) >= 2 {
		return fmt.Errorf("execution error at (%s): %s", string(location), parts[1])
	}

	return err
}

func allTemplates(c *chart.Chart, vals chartutil.Values) map[string]renderable {
	templates := make(map[string]renderable)

	next := map[string]interface{}{
		"Chart":        c.Metadata,
		"Files":        newFiles(c.Files),
		"Release":      vals["Release"],
		"Capabilities": vals["Capabilities"],
		"Values":       vals["Values"],
	}

	chartPath := c.ChartPath()
	for _, t := range c.Templates {
		templates[path.Join(chartPath, t.Name)] = renderable{
			tpl:      string(t.Data),
			vals:     next,
			basePath: path.Join(chartPath, "templates"),
		}
	}

	return templates
}

func sortTemplates(tpls map[string]renderable) []string {
	keys := make([]string, len(tpls))
	i := 0
	for key := range tpls {
		keys[i] = key
		i++
	}
	sort.Sort(sort.Reverse(byPathLen(keys)))
	return keys
}

type byPathLen []string

func (p byPathLen) Len() int      { return len(p) }
func (p byPathLen) Swap(i, j int) { p[j], p[i] = p[i], p[j] }
func (p byPathLen) Less(i, j int) bool {
	a, b := p[i], p[j]
	ca, cb := strings.Count(a, "/"), strings.Count(b, "/")
	if ca == cb {
		return strings.Compare(a, b) == -1
	}
	return ca < cb
}
