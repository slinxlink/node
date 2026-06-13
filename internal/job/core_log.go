package job

import (
	"os"
	"strings"
	"time"

	"github.com/slinxlink/node/internal/database"
)

func CoreLogRotate() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))
			var c database.Core
			database.DB.First(&c)
			trimCoreLog(c.LogPath)
		}
	}()
}

func trimCoreLog(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) <= 10000 {
		return
	}
	lines = lines[len(lines)-10000:]
	os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
