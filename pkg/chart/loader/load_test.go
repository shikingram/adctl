package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadChart(t *testing.T) {
	chart, err := LoadChart("testdata/chart")

	assert.Nil(t, err)
	assert.NotNil(t, chart)

	assert.NotNil(t, chart.Values["hello"])
	assert.Equal(t, 5, len(chart.Templates))
	assert.Equal(t, 7, len(chart.Raw))
}
