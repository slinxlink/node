package api

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/slinxlink/node/internal/database"
)

func tailLines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

func parseSlinxLog(line string) logLine {
	// 2026-06-08 19:44:50 INFO board:SLINX 拉取到 3 个用户
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return logLine{Level: "INFO", Message: line}
	}
	return logLine{
		Time:    parts[0] + " " + parts[1],
		Level:   parts[2],
		Message: parts[3],
	}
}
func parseCoreLog(line string) logLine {
	// +0800 2026-06-08 17:01:07 INFO inbound/hysteria2[hysteria-8998]: ...
	parts := strings.SplitN(line, " ", 5)
	if len(parts) < 5 {
		return logLine{Level: "INFO", Message: line}
	}
	return logLine{
		Time:    parts[1] + " " + parts[2],
		Level:   parts[3],
		Message: parts[4],
	}
}

func logWS(c *gin.Context, pathFn func() string, parseFn func(string) logLine) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	path := pathFn()

	lines, err := tailLines(path, 500)
	if err == nil {
		for _, line := range lines {
			if line == "" {
				continue
			}
			data, _ := json.Marshal(parseFn(line))
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	f.Seek(0, 2) // 跳到末尾

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			// EOF，等一会儿再试
			time.Sleep(500 * time.Millisecond)
			continue
		}
		line = strings.TrimRight(line, "\n")
		if line == "" {
			continue
		}
		data, _ := json.Marshal(parseFn(line))
		if conn.WriteMessage(websocket.TextMessage, data) != nil {
			return
		}
	}
}

func SlinxLog(c *gin.Context) {
	logWS(c, func() string {
		var config database.Config
		database.DB.First(&config)
		return config.LogPath
	}, parseSlinxLog)
}

func CoreLog(c *gin.Context) {
	logWS(c, func() string {
		var core database.Core
		database.DB.First(&core)
		return core.LogPath
	}, parseCoreLog)
}
