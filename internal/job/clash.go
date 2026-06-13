package job

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/slinxlink/node/internal/util"
)

const ClashTemplatePath = "data/clash_template.yaml"

var (
	baseLoyal = "https://raw.githubusercontent.com/Loyalsoldier/clash-rules/release/"
	baseACL   = "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/"
)

type ruleEntry struct {
	url   string
	group string
}

var rulesSources = []ruleEntry{
	// 🤖 AI
	{baseACL + "Ruleset/AI.list", "🤖 AI"},
	{baseACL + "Ruleset/OpenAi.list", "🤖 AI"},
	{baseACL + "Ruleset/Claude.list", "🤖 AI"},
	{baseACL + "Ruleset/ClaudeAI.list", "🤖 AI"},
	{baseACL + "Ruleset/Gemini.list", "🤖 AI"},

	// 🎬 媒体
	{baseACL + "Ruleset/Netflix.list", "🎬 媒体"},
	{baseACL + "Ruleset/NetflixIP.list", "🎬 媒体"},
	{baseACL + "Ruleset/DisneyPlus.list", "🎬 媒体"},
	{baseACL + "Ruleset/YouTube.list", "🎬 媒体"},
	{baseACL + "Ruleset/YouTubeMusic.list", "🎬 媒体"},
	{baseACL + "Ruleset/AppleTV.list", "🎬 媒体"},
	{baseACL + "Ruleset/AppleNews.list", "🎬 媒体"},
	{baseACL + "Ruleset/Spotify.list", "🎬 媒体"},
	{baseACL + "Ruleset/HBO.list", "🎬 媒体"},
	{baseACL + "Ruleset/Hulu.list", "🎬 媒体"},
	{baseACL + "Ruleset/TikTok.list", "🎬 媒体"},
	{baseACL + "Ruleset/Twitch.list", "🎬 媒体"},
	{baseACL + "Ruleset/AbemaTV.list", "🎬 媒体"},
	{baseACL + "Ruleset/DAZN.list", "🎬 媒体"},
	{baseACL + "Ruleset/Niconico.list", "🎬 媒体"},
	{baseACL + "Ruleset/Pixiv.list", "🎬 媒体"},
	{baseACL + "Ruleset/Pandora.list", "🎬 媒体"},
	{baseACL + "Ruleset/BBC.list", "🎬 媒体"},
	{baseACL + "Ruleset/F1.list", "🎬 媒体"},

	// 💬 通讯
	{baseLoyal + "telegramcidr.txt", "💬 通讯"},
	{baseACL + "Telegram.list", "💬 通讯"},
	{baseACL + "Ruleset/Discord.list", "💬 通讯"},
	{baseACL + "Ruleset/Twitter.list", "💬 通讯"},
	{baseACL + "Ruleset/Instagram.list", "💬 通讯"},
	{baseACL + "Ruleset/Facebook.list", "💬 通讯"},
	{baseACL + "Ruleset/Line.list", "💬 通讯"},
	{baseACL + "Ruleset/Whatsapp.list", "💬 通讯"},

	// 🎮 游戏
	{baseACL + "Ruleset/Steam.list", "🎮 游戏"},
	{baseACL + "Ruleset/Epic.list", "🎮 游戏"},
	{baseACL + "Ruleset/Blizzard.list", "🎮 游戏"},
	{baseACL + "Ruleset/Nintendo.list", "🎮 游戏"},
	{baseACL + "Ruleset/Xbox.list", "🎮 游戏"},
	{baseACL + "Ruleset/Origin.list", "🎮 游戏"},
	{baseACL + "Ruleset/Sony.list", "🎮 游戏"},

	// 🍎 苹果
	{baseLoyal + "apple.txt", "🍎 苹果"},
	{baseLoyal + "icloud.txt", "🍎 苹果"},
	{baseACL + "Apple.list", "🍎 苹果"},

	// 🔍 谷歌
	{baseACL + "Ruleset/Google.list", "🔍 谷歌"},
	{baseACL + "Ruleset/GoogleFCM.list", "🔍 谷歌"},

	// Ⓜ️ 微软
	{baseACL + "Microsoft.list", "Ⓜ️ 微软"},
	{baseACL + "OneDrive.list", "Ⓜ️ 微软"},
	{baseACL + "Bing.list", "Ⓜ️ 微软"},
}

func ClashTemplateRefresh() {
	go func() {
		if _, err := os.Stat(ClashTemplatePath); os.IsNotExist(err) {
			generateClashTemplate()
		}
		for {
			now := time.Now()
			daysUntilMonday := (time.Monday - now.Weekday() + 7) % 7
			if daysUntilMonday == 0 {
				daysUntilMonday = 7
			}
			next := time.Date(now.Year(), now.Month(), now.Day()+int(daysUntilMonday), 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))
			generateClashTemplate()
		}
	}()
}

func generateClashTemplate() {
	rules, err := fetchAllRules()
	if err != nil {
		util.Error("[clash] 规则拉取失败: %v", err)
		return
	}

	if err := writeTemplate(rules); err != nil {
		util.Error("[clash] 模板写入失败: %v", err)
		return
	}

	util.Info("[clash] 模板已更新，共 %d 条规则", len(rules))
}

func fetchAllRules() ([]string, error) {
	seen := make(map[string]bool)
	var rules []string

	client := &http.Client{Timeout: 15 * time.Second}

	for _, entry := range rulesSources {
		resp, err := client.Get(entry.url)
		if err != nil {
			util.Warn("[clash] 跳过 %s: %v", entry.url, err)
			continue
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				continue
			}
			if line == "payload:" || line == "---" {
				continue
			}
			line = strings.TrimPrefix(line, "- ")
			line = strings.Trim(line, `'"`)

			var rule string
			switch {
			case strings.HasPrefix(line, "IP-CIDR,"),
				strings.HasPrefix(line, "IP-CIDR6,"):
				line = strings.TrimSuffix(line, ",no-resolve")
				rule = line + "," + entry.group + ",no-resolve"
			case strings.HasPrefix(line, "DOMAIN,"),
				strings.HasPrefix(line, "DOMAIN-SUFFIX,"),
				strings.HasPrefix(line, "DOMAIN-KEYWORD,"),
				strings.HasPrefix(line, "PROCESS-NAME,"):
				rule = line + "," + entry.group
			case strings.HasPrefix(line, "+.") || strings.HasPrefix(line, "*."):
				rule = "DOMAIN-SUFFIX," + line[2:] + "," + entry.group
			default:
				rule = "DOMAIN-SUFFIX," + line + "," + entry.group
			}

			var key string
			if strings.HasPrefix(line, "IP-CIDR,") || strings.HasPrefix(line, "IP-CIDR6,") {
				// 去掉 ,group,no-resolve 两段
				idx := strings.LastIndex(rule, ",")
				key = rule[:strings.LastIndex(rule[:idx], ",")]
			} else {
				key = rule[:strings.LastIndex(rule, ",")]
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			rules = append(rules, rule)
		}
		resp.Body.Close()
		if err := scanner.Err(); err != nil {
			util.Warn("[clash] 读取 %s 失败: %v", entry.url, err)
		}
	}

	rules = append(rules,
		"GEOIP,LAN,🏠 局域网",
		"GEOSITE,CN,🇨🇳 大陆",
		"GEOIP,CN,🇨🇳 大陆",
		"MATCH,🐟 漏网之鱼",
	)

	return rules, nil
}

func writeTemplate(rules []string) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`# Clash 订阅模板 - 由 slinx-node 自动生成，请勿手动修改
mixed-port: 7890
allow-lan: true
mode: Rule
log-level: error
external-controller: 0.0.0.0:9090

proxies: []

proxy-groups:
  - name: 🚀 代理
    type: select
    proxies:
# __PROXIES__
  - name: 🤖 AI
    type: select
    proxies:
      - 🚀 代理
      - 🎯 直连
# __PROXIES__
  - name: 🎬 媒体
    type: select
    proxies:
      - 🚀 代理
      - 🎯 直连
# __PROXIES__
  - name: 💬 通讯
    type: select
    proxies:
      - 🚀 代理
      - 🎯 直连
# __PROXIES__
  - name: 🎮 游戏
    type: select
    proxies:
      - 🚀 代理
      - 🎯 直连
# __PROXIES__
  - name: 🍎 苹果
    type: select
    proxies:
      - 🎯 直连
      - 🚀 代理
# __PROXIES__
  - name: 🔍 谷歌
    type: select
    proxies:
      - 🚀 代理
      - 🎯 直连
# __PROXIES__
  - name: Ⓜ️ 微软
    type: select
    proxies:
      - 🎯 直连
      - 🚀 代理
# __PROXIES__
  - name: 🏠 局域网
    type: select
    proxies:
      - 🎯 直连
      - 🚀 代理
# __PROXIES__
  - name: 🇨🇳 大陆
    type: select
    proxies:
      - 🎯 直连
      - 🚀 代理
# __PROXIES__
  - name: 🐟 漏网之鱼
    type: select
    proxies:
      - 🚀 代理
      - 🎯 直连
# __PROXIES__
  - name: 🎯 直连
    type: select
    proxies:
      - DIRECT
  - name: 🛑 阻断
    type: select
    proxies:
      - REJECT

rules:
`)

	for _, rule := range rules {
		sb.WriteString("  - ")
		sb.WriteString(rule)
		sb.WriteString("\n")
	}

	return os.WriteFile(ClashTemplatePath, []byte(sb.String()), 0644)
}
