package chartutil

import (
	"fmt"
	"log"

	"github.com/shikingram/adctl/pkg/chart"

	"github.com/mitchellh/copystructure"
)

func CoalesceValues(chrt *chart.Chart, vals map[string]interface{}) (Values, error) {
	v, err := copystructure.Copy(vals)
	if err != nil {
		return vals, err
	}

	valsCopy := v.(map[string]interface{})
	// if we have an empty map, make sure it is initialized
	if valsCopy == nil {
		valsCopy = make(map[string]interface{})
	}
	return coalesce(log.Printf, chrt, valsCopy, "")
}

type printFn func(format string, v ...interface{})

func coalesce(printf printFn, ch *chart.Chart, dest map[string]interface{}, prefix string) (map[string]interface{}, error) {
	coalesceValues(printf, ch, dest, prefix)
	return dest, nil
}

// coalesceValues builds up a values map for a particular chart.
//
// Values in v will override the values in the chart.
func coalesceValues(printf printFn, c *chart.Chart, v map[string]interface{}, prefix string) {
	subPrefix := concatPrefix(prefix, c.Metadata.Name)
	for key, val := range c.Values {
		if value, ok := v[key]; ok {
			if value == nil {
				delete(v, key)
			} else if dest, ok := value.(map[string]interface{}); ok {
				// if v[key] is a table, merge nv's val table into v[key].
				src, ok := val.(map[string]interface{})
				if !ok {
					// If the original value is nil, there is nothing to coalesce, so we don't print
					// the warning
					if val != nil {
						printf("warning: skipped value for %s.%s: Not a table.", subPrefix, key)
					}
				} else {
					// Because v has higher precedence than nv, dest values override src
					// values.
					coalesceTablesFullKey(printf, dest, src, concatPrefix(subPrefix, key))
				}
			}
		} else {
			// If the key is not in v, copy it from nv.
			v[key] = val
		}
	}
}

func concatPrefix(a, b string) string {
	if a == "" {
		return b
	}
	return fmt.Sprintf("%s.%s", a, b)
}

func coalesceTablesFullKey(printf printFn, dst, src map[string]interface{}, prefix string) map[string]interface{} {
	// When --reuse-values is set but there are no modifications yet, return new values
	if src == nil {
		return dst
	}
	if dst == nil {
		return src
	}
	// Because dest has higher precedence than src, dest values override src
	// values.
	for key, val := range src {
		fullkey := concatPrefix(prefix, key)
		if dv, ok := dst[key]; ok && dv == nil {
			delete(dst, key)
		} else if !ok {
			dst[key] = val
		} else if istable(val) {
			if istable(dv) {
				coalesceTablesFullKey(printf, dv.(map[string]interface{}), val.(map[string]interface{}), fullkey)
			} else {
				printf("warning: cannot overwrite table with non table for %s (%v)", fullkey, val)
			}
		} else if istable(dv) && val != nil {
			printf("warning: destination for %s is a table. Ignoring non-table value (%v)", fullkey, val)
		}
	}
	return dst
}

func istable(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}
