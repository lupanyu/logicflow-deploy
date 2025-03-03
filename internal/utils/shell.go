package utils

import (
	"os/exec"
	"time"
)

func RunShell(cmd string) ([]byte, error) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return out, err
}

func GetNowTime() *time.Time {
	res := time.Now()
	return &res
}
