package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/slinxlink/node/internal/database"
)

type Connection struct {
	ID       string `json:"id"`
	Metadata struct {
		Network  string   `json:"network"`
		SourceIP string   `json:"sourceIP"`
		Chains   []string `json:"chains"`
	} `json:"metadata"`
}

type ConnectionsResponse struct {
	Connections []Connection `json:"connections"`
}

type OnlineUser struct {
	Name string
	IP   string
}

type connEntry struct {
	IP       string
	UserName string
}

func GetOnlineUsers() ([]OnlineUser, error) {
	res, err := http.Get("http://127.0.0.1:9090/connections")
	if err != nil {
		return nil, fmt.Errorf("获取连接失败: %w", err)
	}
	defer res.Body.Close()

	var resp ConnectionsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("解析连接失败: %w", err)
	}

	activeIPs := make(map[string]bool)
	for _, conn := range resp.Connections {
		if conn.Metadata.SourceIP != "" {
			activeIPs[conn.Metadata.SourceIP] = true
		}
	}

	if len(activeIPs) == 0 {
		return nil, nil
	}

	lines, err := readRecentLog()
	if err != nil {
		return nil, err
	}

	connMap := make(map[string]*connEntry)
	for _, line := range lines {
		parseLogLineInto(line, connMap)
	}

	seen := make(map[string]map[string]bool)
	for _, e := range connMap {
		if e.IP == "" || e.UserName == "" {
			continue
		}
		if !activeIPs[e.IP] {
			continue
		}
		if seen[e.UserName] == nil {
			seen[e.UserName] = make(map[string]bool)
		}
		seen[e.UserName][e.IP] = true
	}

	var result []OnlineUser
	for name, ips := range seen {
		for ip := range ips {
			result = append(result, OnlineUser{Name: name, IP: ip})
		}
	}
	return result, nil
}

func readRecentLog() ([]string, error) {
	var c database.Core
	database.DB.First(&c)
	if !c.LogEnable || c.LogPath == "" {
		return nil, nil
	}

	f, err := os.Open(c.LogPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := stat.Size()
	if size == 0 {
		return nil, nil
	}

	readSize := int64(512 * 1024)
	if size < readSize {
		readSize = size
	}

	f.Seek(-readSize, io.SeekEnd)
	buf := make([]byte, readSize)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buf = buf[:n]

	cutoff := time.Now().Add(-1 * time.Minute)
	var lines []string
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}
		t, ok := parseLogTime(line)
		if !ok {
			continue
		}
		if t.After(cutoff) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func parseLogTime(line string) (time.Time, bool) {
	// 格式: +0800 2026-06-27 13:04:37 INFO ...
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return time.Time{}, false
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", parts[1]+" "+parts[2], time.Local)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func parseLogLineInto(line string, connMap map[string]*connEntry) {
	idx := strings.Index(line, " [")
	if idx == -1 {
		return
	}
	rest := line[idx+2:]
	spaceIdx := strings.Index(rest, " ")
	if spaceIdx == -1 {
		return
	}
	connID := rest[:spaceIdx]

	if connID == "" {
		return
	}
	for _, c := range connID {
		if c < '0' || c > '9' {
			return
		}
	}

	if strings.Contains(line, "inbound connection from ") || strings.Contains(line, "inbound packet connection from ") {
		idx := strings.Index(line, " from ")
		if idx == -1 {
			return
		}
		addr := strings.TrimSpace(line[idx+len(" from "):])
		ip := addr[:strings.LastIndex(addr, ":")]
		ip = strings.TrimPrefix(ip, "[")
		ip = strings.TrimSuffix(ip, "]")
		if connMap[connID] == nil {
			connMap[connID] = &connEntry{}
		}
		connMap[connID].IP = ip

	} else if strings.Contains(line, "] inbound connection to ") || strings.Contains(line, "] inbound packet connection to ") {
		s := strings.Index(line, "[BoardUser_")
		if s == -1 {
			return
		}
		e := strings.Index(line[s:], "]")
		if e == -1 {
			return
		}
		userName := line[s+1 : s+e]
		if connMap[connID] == nil {
			connMap[connID] = &connEntry{}
		}
		connMap[connID].UserName = userName
	}
}
