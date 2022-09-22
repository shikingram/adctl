package deploy

import (
	"github.com/bitfield/script"
)

func run(cmd string) error {
	_, err := script.Exec(cmd).Stdout()
	if err != nil {
		return err
	}
	return nil
}
