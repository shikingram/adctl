package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
)

const ps = `docker ps --filter "label=com.docker.compose.project" -q`
const inspect = `xargs docker inspect --format='{{index .Config.Labels "com.docker.compose.project"}}'`

// start docker-compose
func Start(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	pname = "ad-" + name + "-" + pname
	cmd := fmt.Sprintf("docker-compose -f %s -p %s up -d --remove-orphans", file, pname)
	return exec(cmd)
}

func exec(cmd string) error {
	fmt.Printf("# exec cmd: %s \n", cmd)
	_, err := script.Exec(cmd).Stdout()
	return err
}

// stop docker-compose
func Stop(file string) error {
	cmd := fmt.Sprintf("docker-compose -f %s stop -t 20", file)
	return exec(cmd)
}

// stop and delete docker container
func Down(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	pname = "ad-" + name + "-" + pname
	cmd := fmt.Sprintf("docker-compose -f %s -p %s down --remove-orphans", file, pname)
	return exec(cmd)
}

func Restart(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	pname = "ad-" + name + "-" + pname
	cmd := fmt.Sprintf("docker-compose -f %s -p %s up -d --remove-orphans && docker-compose -f %s -p %s restart", file, pname, file, pname)
	return exec(cmd)
}

func CheckReleaseDeploy(name string) (int, error) {

	if l, err := script.Exec(ps).CountLines(); err == nil && l > 0 {
		return script.Exec(ps).
			Exec(inspect).
			Exec(`sort`).
			Exec(`uniq`).
			Match(name).
			CountLines()
	}

	return 0, nil

}

func ListRelease() error {

	if l, err := script.Exec(ps).CountLines(); err == nil && l > 0 {
		_, e := script.Exec(ps).
			Exec(inspect).
			Exec(`sort`).
			Exec(`uniq`).
			Stdout()
		if e != nil {
			return e
		}
	}

	return nil
}
