package deploy

import (
	"fmt"

	"github.com/bitfield/script"
)

func CreateNetWork(name string) error {
	numErrors, err := script.Exec(fmt.Sprintf("docker network inspect %s", name)).Match("Error").CountLines()
	if err != nil {
		return err
	}
	if numErrors > 0 {
		return run(fmt.Sprintf("docker network create %s ", name))
	} else {
		fmt.Printf("network %s is already exists \n", name)
	}
	return err
}
