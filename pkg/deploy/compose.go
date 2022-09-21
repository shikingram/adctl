package deploy

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// start docker-compose
func Start(file string) error {
	pName := file[strings.LastIndex(file, "/")+8 : strings.LastIndex(file, ".")]
	return run("bash", fmt.Sprintf("docker-compose -f %s up -d -p %s --remove-orphans", file, pName))
}

// stop docker-compose
func Stop(file string) error {
	return run("bash", fmt.Sprintf("docker-compose -f %s stop -t 20", file))
}

// stop and delete docker container
func Clean(file string) error {
	if err := deleteDir(path.Base(file)); err != nil {
		return err
	}
	if err := run("bash", fmt.Sprintf("docker-compose -f %s down --remove-orphans", file)); err != nil {
		return err
	}
	return nil
}

func deleteDir(pathList ...string) (err error) {
	for _, file := range pathList {
		if err = os.RemoveAll(file); err != nil {
			return
		}
	}
	return
}
