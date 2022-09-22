package deploy

import (
	"fmt"
	"testing"
)

func TestCheckReleaseDeploy(t *testing.T) {
	i, e := CheckReleaseDeploy("mysql")
	if e == nil {
		fmt.Println(i)
	}
}
