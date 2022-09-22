package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
)

// start docker-compose
func Start(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	_, err := script.Exec(fmt.Sprintf("docker-compose -f %s -p ad-%s-%s up -d --remove-orphans", file, name, pname)).Stdout()
	return err
}

// stop docker-compose
func Stop(file string) error {
	_, err := script.Exec(fmt.Sprintf("docker-compose -f %s stop -t 20", file)).Stdout()
	return err
}

// stop and delete docker container
func Down(file, name string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".yaml")]
	_, err := script.Exec(fmt.Sprintf("docker-compose -f %s -p ad-%s-%s down --remove-orphans", file, name, pname)).Stdout()
	return err
}

func CheckReleaseDeploy(name string) (int, error) {
	return script.Exec(`docker ps --filter "label=com.docker.compose.project" -q`).
		Exec(`xargs docker inspect --format='{{index .Config.Labels "com.docker.compose.project"}}'`).
		Exec(`sort`).
		Exec(`uniq`).
		Match(name).
		CountLines()
}
