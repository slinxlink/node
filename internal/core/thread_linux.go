//go:build linux

package core

import (
	"fmt"
	"os"
	"strings"
)

func getThread(pid int) int {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "Threads:") {
			threads := 0
			fmt.Sscanf(strings.TrimPrefix(line, "Threads:"), "%d", &threads)
			return threads
		}
	}
	return 0
}
