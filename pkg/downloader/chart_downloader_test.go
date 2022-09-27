package downloader

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLParse(t *testing.T) {
	u, err := url.Parse("https://charts.kubevela.net/core")
	assert.Nil(t, err)
	fmt.Printf("u:%+#v \n", u)
}
