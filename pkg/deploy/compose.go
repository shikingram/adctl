package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// start docker-compose
func Start(file string) error {
	basefile := filepath.Base(file)
	pname := basefile[7:strings.LastIndex(basefile, ".")]
	return run(fmt.Sprintf("docker-compose -f %s -p %s up -d --remove-orphans", file, pname))
}

// stop docker-compose
func Stop(file string) error {
	return run(fmt.Sprintf("docker-compose -f %s stop -t 20", file))
}

// stop and delete docker container
func Down(file string) error {
	os.RemoveAll(file)
	return run(fmt.Sprintf("docker-compose -f %s down --remove-orphans", file))
}
