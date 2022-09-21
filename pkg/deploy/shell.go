package deploy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

func run(shell, cmd string) error {
	exe := exec.Command(shell, "-c", cmd)
	stdout, err := exe.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := exe.StderrPipe()
	if err != nil {
		return err
	}
	err = exe.Start()
	if err != nil {
		return err
	}

	go outToLog(stdout)
	res, _ := ioutil.ReadAll(stderr)
	stderrString := string(res)
	err = exe.Wait()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s;%s", err.Error(), stderrString))
	}
	return nil
}

func outToLog(output io.ReadCloser) {
	var tmp [1024]byte
	for {
		n, err := output.Read(tmp[:])
		if err != nil {
			break
		}
		fmt.Printf("%s \n", strings.Replace(string(tmp[:n]), "\n", " ", -1))
	}
}
