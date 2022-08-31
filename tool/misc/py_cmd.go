package misc

import (
	"os/exec"
)

func CmdPy(pythonPath string, cmdArgs ...string) (result string, err error) {
	args := []string{pythonPath}
	args = append(args, cmdArgs...)
	out, err := exec.Command("python", args...).Output()
	if err != nil {
		return "", err
	}
	re := string(out)
	return re, nil
}
