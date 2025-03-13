package utils

import (
	"os/exec"
)

func RunShell(cmd string) ([]byte, error) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return out, err
}
