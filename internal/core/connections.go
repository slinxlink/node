package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	seen := make(map[string]map[string]bool) // name -> set of IPs
	for _, conn := range resp.Connections {
		// Chains 里面包含用户名，格式类似 ["BoardUser_1", "DIRECT"]
		for _, chain := range conn.Metadata.Chains {
			if strings.HasPrefix(chain, "BoardUser_") {
				if seen[chain] == nil {
					seen[chain] = make(map[string]bool)
				}
				seen[chain][conn.Metadata.SourceIP] = true
			}
		}
	}

	var result []OnlineUser
	for name, ips := range seen {
		for ip := range ips {
			result = append(result, OnlineUser{Name: name, IP: ip})
		}
	}
	return result, nil
}
