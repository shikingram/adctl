package deploy

import "fmt"

func CreateNetWork(name string) error {
	return run("bash", fmt.Sprintf("docker network inspect %s > /dev/null || docker network create  %s ", name, name))
}
