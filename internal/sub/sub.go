package sub

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/job"
	"github.com/slinxlink/node/internal/util"
)

// ── 订阅入口 ─────────────────────────────────────────────────────────

func Sub(token string) string {
	user, inbounds := getUser(token)
	if user == nil {
		return ""
	}

	host := getHost()
	var uris []string
	for _, inbound := range inbounds {
		uri := dispatch(*user, inbound, host)
		if uri != "" {
			uris = append(uris, uri)
		}
	}
	return strings.Join(uris, "\n")
}

func Clash(token string) string {
	user, inbounds := getUser(token)
	if user == nil {
		return ""
	}

	host := getHost()
	var proxies []string
	for _, inbound := range inbounds {
		proxy := clashProxy(*user, inbound, host)
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}
	return renderClash(proxies)
}

// ── 信息查询 ─────────────────────────────────────────────────────────

type Data struct {
	User     database.User      `json:"user"`
	Inbounds []database.Inbound `json:"inbounds"`
	Uris     []string           `json:"uris"`
	Urls     []string           `json:"urls"`
}

func Info(token string) *Data {
	var user database.User
	if database.DB.Where("token = ? AND enable = ?", token, true).First(&user).Error != nil {
		return nil
	}

	var ids []int
	if err := json.Unmarshal([]byte(user.Inbounds), &ids); err != nil {
		return nil
	}

	var inbounds []database.Inbound
	database.DB.Where("id IN ? AND enable = ?", ids, true).Find(&inbounds)

	host := getHost()
	var uris []string
	for _, inbound := range inbounds {
		uri := dispatch(user, inbound, host)
		if uri != "" {
			uris = append(uris, uri)
		}
	}

	return &Data{
		User:     user,
		Inbounds: inbounds,
		Uris:     uris,
		Urls:     Url(token),
	}
}

func Url(token string) []string {
	var cfg database.Config
	database.DB.First(&cfg)

	host := cfg.Domain
	if host == "" {
		host = cfg.IPv4
	}

	scheme := "http"
	if cfg.Domain != "" {
		scheme = "https"
	}

	base := fmt.Sprintf("%s://%s:%d", scheme, host, cfg.SubPort)

	return []string{
		fmt.Sprintf("%s%s/%s", base, cfg.SubPath, token),
		fmt.Sprintf("%s%s/%s", base, cfg.ClashPath, token),
	}
}

func Uri(user database.User, inbound database.Inbound) string {
	host := getHost()
	return dispatch(user, inbound, host)
}

// ── 协议分发 ─────────────────────────────────────────────────────────

func dispatch(user database.User, inbound database.Inbound, host string) string {
	switch inbound.Protocol {
	case "vless":
		return vless(user.UUID, host, inbound)
	case "vmess":
		return vmess(user.UUID, host, inbound)
	case "hysteria":
		return hysteria(user.Password, host, inbound)
	default:
		return ""
	}
}

func clashProxy(user database.User, inbound database.Inbound, host string) string {
	switch inbound.Protocol {
	case "vless":
		return vlessClash(user.UUID, host, inbound)
	case "vmess":
		return vmessClash(user.UUID, host, inbound)
	case "hysteria":
		return hysteriaClash(user.Password, host, inbound)
	default:
		return ""
	}
}

// ── 内部工具 ─────────────────────────────────────────────────────────

func getHost() string {
	var config database.Config
	database.DB.First(&config)
	if config.Domain != "" {
		return config.Domain
	}
	return config.IPv4
}

func getUser(token string) (*database.User, []database.Inbound) {
	var user database.User
	if database.DB.Where("token = ? AND enable = ?", token, true).First(&user).Error != nil {
		return nil, nil
	}
	var ids []int
	if err := json.Unmarshal([]byte(user.Inbounds), &ids); err != nil {
		return nil, nil
	}
	var inbounds []database.Inbound
	database.DB.Where("id IN ? AND enable = ?", ids, true).Find(&inbounds)
	return &user, inbounds
}

func renderClash(proxies []string) string {
	data, err := os.ReadFile(job.ClashTemplatePath)
	if err != nil {
		util.Error("[clash] 读取模板失败: %v", err)
		return ""
	}

	var names []string
	for _, p := range proxies {
		for _, line := range strings.Split(p, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- name: ") {
				names = append(names, strings.TrimPrefix(trimmed, "- name: "))
				break
			}
		}
	}

	nameBlock := ""
	for _, name := range names {
		nameBlock += "      - " + name + "\n"
	}

	yaml := string(data)
	proxyBlock := strings.Join(proxies, "")
	yaml = strings.Replace(yaml, "proxies: []", "proxies:\n"+proxyBlock, 1)
	yaml = strings.ReplaceAll(yaml, "# __PROXIES__\n", nameBlock)

	return yaml
}
