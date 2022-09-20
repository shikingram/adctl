package chart

import "fmt"

// ValidationError represents a data validation error.
type ValidationError string

func (v ValidationError) Error() string {
	return "validation: " + string(v)
}

// ValidationErrorf takes a message and formatting options and creates a ValidationError
func ValidationErrorf(msg string, args ...interface{}) ValidationError {
	return ValidationError(fmt.Sprintf(msg, args...))
}
