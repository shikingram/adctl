package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	err := run("echo 127.0.0.1")
	assert.Nil(t, err)
}
