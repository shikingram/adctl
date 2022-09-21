package chartutil

import (
	"fmt"
)

// ErrNoTable indicates that a chart does not have a matching table.
type ErrNoTable struct {
	Key string
}

func (e ErrNoTable) Error() string { return fmt.Sprintf("%q is not a table", e.Key) }

// ErrNoValue indicates that Values does not contain a key with a value
type ErrNoValue struct {
	Key string
}

func (e ErrNoValue) Error() string { return fmt.Sprintf("%q is not a value", e.Key) }
