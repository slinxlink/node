package sub

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/slinxlink/node/internal/config"
	"github.com/slinxlink/node/internal/database"
	tpl "github.com/slinxlink/node/internal/sub/template"
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

func Clash(token string) (string, string) {
	user, inbounds := getUser(token)
	if user == nil {
		return "", ""
	}

	host := getHost()
	var proxies []string
	for _, inbound := range inbounds {
		proxy := dispatchClash(*user, inbound, host)
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}

	name := util.SanitizeFileName(user.Name)
	if name == "" {
		name = "SLINX"
	}
	return tpl.RenderClash(proxies), name
}

func Surge(token string) (string, string) {
	user, inbounds := getUser(token)
	if user == nil {
		return "", ""
	}

	host := getHost()
	var proxies []string
	var names []string
	for _, inbound := range inbounds {
		proxy := dispatchSurge(*user, inbound, host)
		if proxy != "" {
			proxies = append(proxies, proxy)
			names = append(names, inbound.Name)
		}
	}

	name := util.SanitizeFileName(user.Name)
	if name == "" {
		name = "SLINX"
	}
	return tpl.RenderSurge(proxies, names), name
}

// ── 信息查询 ─────────────────────────────────────────────────────────

type Data struct {
	User     database.User      `json:"user"`
	Inbounds []database.Inbound `json:"inbounds"`
	Uris     []string           `json:"uris"`
	Urls     []string           `json:"urls"`
	Jsons    []string           `json:"jsons"`
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
	var jsons []string
	for _, inbound := range inbounds {
		uri := dispatch(user, inbound, host)
		if uri != "" {
			uris = append(uris, uri)
		}
		jsons = append(jsons, Json(user, inbound, "singbox"))
	}

	return &Data{
		User:     user,
		Inbounds: inbounds,
		Uris:     uris,
		Urls:     Url(token),
		Jsons:    jsons,
	}
}

func Url(token string) []string {
	cfg := config.Config

	host := cfg.Domain
	if host == "" {
		host = cfg.IPv4
	}

	scheme := "http"
	if cfg.Domain != "" {
		scheme = "https"
	}

	base := fmt.Sprintf("%s://%s:%d", scheme, host, cfg.SubPort)
	sub := fmt.Sprintf("%s%s/%s", base, cfg.SubPath, token)

	return []string{
		sub,
		sub + "/clash",
		sub + "/surge",
	}
}

func Uri(user database.User, inbound database.Inbound) string {
	host := getHost()
	return dispatch(user, inbound, host)
}

func Json(user database.User, inbound database.Inbound, format string) string {
	host := getHost()
	switch format {
	case "singbox":
		outbound := dispatchSingBox(user, inbound, host)
		if outbound == "" {
			return ""
		}
		return tpl.RenderSingBox([]string{outbound})
	default:
		return ""
	}
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
	case "trojan":
		return trojan(user.Password, host, inbound)
	default:
		return ""
	}
}

func dispatchClash(user database.User, inbound database.Inbound, host string) string {
	switch inbound.Protocol {
	case "vless":
		return vlessClash(user.UUID, host, inbound)
	case "vmess":
		return vmessClash(user.UUID, host, inbound)
	case "hysteria":
		return hysteriaClash(user.Password, host, inbound)
	case "trojan":
		return trojanClash(user.Password, host, inbound)
	default:
		return ""
	}
}

func dispatchSurge(user database.User, inbound database.Inbound, host string) string {
	switch inbound.Protocol {
	case "vmess":
		return vmessSurge(user.UUID, host, inbound)
	case "hysteria":
		return hysteriaSurge(user.Password, host, inbound)
	case "trojan":
		return trojanSurge(user.Password, host, inbound)
	default:
		return ""
	}
}

func dispatchSingBox(user database.User, inbound database.Inbound, host string) string {
	switch inbound.Protocol {
	case "vless":
		return vlessSingBox(user.UUID, host, inbound)
	case "vmess":
		return vmessSingBox(user.UUID, host, inbound)
	case "hysteria":
		return hysteriaSingBox(user.Password, host, inbound)
	case "trojan":
		return trojanSingBox(user.Password, host, inbound)
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
