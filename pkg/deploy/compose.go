package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
)

const dprefix = `ad-`

// start docker-compose
func Start(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	pname = dprefix + name + "-" + pname
	cmd := fmt.Sprintf("docker compose -f %s -p %s up -d --remove-orphans", file, pname)
	return exec(cmd)
}

func exec(cmd string) error {
	// fmt.Printf("# exec cmd: %s \n", cmd)
	_, err := script.Exec(cmd).Stdout()
	return err
}

// stop docker-compose
func Stop(file string) error {
	cmd := fmt.Sprintf("docker compose -f %s stop -t 20", file)
	return exec(cmd)
}

// stop and delete docker container
func Down(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	pname = dprefix + name + "-" + pname
	cmd := fmt.Sprintf("docker compose -f %s -p %s down --remove-orphans", file, pname)
	return exec(cmd)
}

func Restart(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	pname = dprefix + name + "-" + pname
	cmd1 := fmt.Sprintf("docker compose -f %s -p %s up -d --remove-orphans", file, pname)

	cmd2 := fmt.Sprintf("docker compose -f %s -p %s restart", file, pname)

	err := exec(cmd1)
	if err != nil {
		return err
	}

	return exec(cmd2)
}

func CheckReleaseDeploy(name string) (int, error) {
	return script.Exec(`docker compose ls`).Match(dprefix).Match(name).CountLines()
}

func ListRelease(name string) error {

	pipe := script.Exec(`docker compose ls`)

	var e error
	if len(name) > 0 {
		_, e = pipe.Match(dprefix).Match(name).Stdout()
	} else {
		_, e = pipe.Match(dprefix).Stdout()
	}
	return e
}
