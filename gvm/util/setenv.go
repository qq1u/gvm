package util

import "os/exec"

func execute(cmd *exec.Cmd) (err error) {
	_, err = cmd.Output()
	return
}
