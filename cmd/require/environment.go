package require

import (
	"errors"

	"github.com/bitfield/script"
)

func Environment() error {
	countNum, err := script.Exec("where docker && where docker-compose").Match("not found").CountLines()
	if err != nil {
		return err
	}
	if countNum > 0 {
		return errors.New("docker or dockercompose not found")
	}
	return nil
}
