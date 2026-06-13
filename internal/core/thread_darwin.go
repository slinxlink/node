//go:build darwin

package core

import (
	"fmt"
	"os/exec"
	"strings"
)

func getThread(pid int) int {
	out, err := exec.Command("ps", "-M", "-p", fmt.Sprintf("%d", pid)).Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 1 {
		return len(lines) - 1
	}
	return 0
}
