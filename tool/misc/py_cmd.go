package misc

import (
	"os/exec"
)

func CmdPy(pythonPath string, cmdArgs ...string) (result string, err error) {
	args := []string{pythonPath}
	args = append(args, cmdArgs...)
	pys := []string{"python", "python3"}
	usePy := ""
	for _, py := range pys {
		_, err = exec.LookPath(py)
		if err == nil {
			usePy = py
			break
		}
	}
	out, err := exec.Command(usePy, args...).Output()
	if err != nil {
		return "", err
	}
	re := string(out)
	return re, nil
}
