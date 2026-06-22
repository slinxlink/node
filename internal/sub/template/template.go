package template

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/slinxlink/node/internal/util"
)

const (
	RulesetPath         = "data/ruleset/template.json"
	ClashTemplatePath   = "data/ruleset/clash.yaml"
	SurgeTemplatePath   = "data/ruleset/surge.conf"
	SingBoxTemplatePath = "data/ruleset/singbox.json"
)

type RuleGroup struct {
	Name  string   `json:"name"`
	Rules []string `json:"rules"`
}

type Ruleset struct {
	Groups []RuleGroup `json:"groups"`
}

var (
	baseLoyal = "https://raw.githubusercontent.com/Loyalsoldier/clash-rules/release/"
	baseACL   = "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/"
)

type ruleEntry struct {
	url   string
	group string
}

var rulesSources = []ruleEntry{
	{baseACL + "Ruleset/AI.list", "AI"},
	{baseACL + "Ruleset/OpenAi.list", "AI"},
	{baseACL + "Ruleset/Claude.list", "AI"},
	{baseACL + "Ruleset/ClaudeAI.list", "AI"},
	{baseACL + "Ruleset/Gemini.list", "AI"},

	{baseACL + "Ruleset/Netflix.list", "MEDIA"},
	{baseACL + "Ruleset/NetflixIP.list", "MEDIA"},
	{baseACL + "Ruleset/DisneyPlus.list", "MEDIA"},
	{baseACL + "Ruleset/YouTube.list", "MEDIA"},
	{baseACL + "Ruleset/YouTubeMusic.list", "MEDIA"},
	{baseACL + "Ruleset/AppleTV.list", "MEDIA"},
	{baseACL + "Ruleset/AppleNews.list", "MEDIA"},
	{baseACL + "Ruleset/Spotify.list", "MEDIA"},
	{baseACL + "Ruleset/HBO.list", "MEDIA"},
	{baseACL + "Ruleset/Hulu.list", "MEDIA"},
	{baseACL + "Ruleset/TikTok.list", "MEDIA"},
	{baseACL + "Ruleset/Twitch.list", "MEDIA"},
	{baseACL + "Ruleset/AbemaTV.list", "MEDIA"},
	{baseACL + "Ruleset/DAZN.list", "MEDIA"},
	{baseACL + "Ruleset/Niconico.list", "MEDIA"},
	{baseACL + "Ruleset/Pixiv.list", "MEDIA"},
	{baseACL + "Ruleset/Pandora.list", "MEDIA"},
	{baseACL + "Ruleset/BBC.list", "MEDIA"},
	{baseACL + "Ruleset/F1.list", "MEDIA"},

	{baseLoyal + "telegramcidr.txt", "SOCIAL"},
	{baseACL + "Telegram.list", "SOCIAL"},
	{baseACL + "Ruleset/Discord.list", "SOCIAL"},
	{baseACL + "Ruleset/Twitter.list", "SOCIAL"},
	{baseACL + "Ruleset/Instagram.list", "SOCIAL"},
	{baseACL + "Ruleset/Facebook.list", "SOCIAL"},
	{baseACL + "Ruleset/Line.list", "SOCIAL"},
	{baseACL + "Ruleset/Whatsapp.list", "SOCIAL"},

	{baseACL + "Ruleset/Steam.list", "GAME"},
	{baseACL + "Ruleset/Epic.list", "GAME"},
	{baseACL + "Ruleset/Blizzard.list", "GAME"},
	{baseACL + "Ruleset/Nintendo.list", "GAME"},
	{baseACL + "Ruleset/Xbox.list", "GAME"},
	{baseACL + "Ruleset/Origin.list", "GAME"},
	{baseACL + "Ruleset/Sony.list", "GAME"},

	{baseLoyal + "apple.txt", "APPLE"},
	{baseLoyal + "icloud.txt", "APPLE"},
	{baseACL + "Apple.list", "APPLE"},

	{baseACL + "Ruleset/Google.list", "GOOGLE"},
	{baseACL + "Ruleset/GoogleFCM.list", "GOOGLE"},

	{baseACL + "Microsoft.list", "MICROSOFT"},
	{baseACL + "OneDrive.list", "MICROSOFT"},
	{baseACL + "Bing.list", "MICROSOFT"},
}

func DownloadRuleset() (int, error) {
	seen := make(map[string]bool)
	groupMap := make(map[string][]string)
	var groupOrder []string

	client := &http.Client{Timeout: 15 * time.Second}

	for _, entry := range rulesSources {
		if _, exists := groupMap[entry.group]; !exists {
			groupOrder = append(groupOrder, entry.group)
			groupMap[entry.group] = []string{}
		}

		resp, err := client.Get(entry.url)
		if err != nil {
			util.Warn("[ruleset] 跳过 %s: %v", entry.url, err)
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
			case strings.HasPrefix(line, "IP-CIDR,"), strings.HasPrefix(line, "IP-CIDR6,"):
				rule = strings.TrimSuffix(line, ",no-resolve")
			case strings.HasPrefix(line, "DOMAIN,"),
				strings.HasPrefix(line, "DOMAIN-SUFFIX,"),
				strings.HasPrefix(line, "DOMAIN-KEYWORD,"),
				strings.HasPrefix(line, "PROCESS-NAME,"):
				rule = line
			case strings.HasPrefix(line, "+.") || strings.HasPrefix(line, "*."):
				rule = "DOMAIN-SUFFIX," + line[2:]
			default:
				rule = "DOMAIN-SUFFIX," + line
			}

			if seen[rule] {
				continue
			}
			seen[rule] = true
			groupMap[entry.group] = append(groupMap[entry.group], rule)
		}
		if err := scanner.Err(); err != nil {
			util.Warn("[ruleset] 读取 %s 失败: %v", entry.url, err)
		}
		resp.Body.Close()
	}

	var groups []RuleGroup
	for _, name := range groupOrder {
		groups = append(groups, RuleGroup{Name: name, Rules: groupMap[name]})
	}
	groups = append(groups, RuleGroup{
		Name:  "_final",
		Rules: []string{"GEOIP,LAN", "GEOSITE,CN", "GEOIP,CN"},
	})

	total := len(seen)

	data, err := json.MarshalIndent(Ruleset{Groups: groups}, "", "  ")
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll("data/ruleset", 0755); err != nil {
		return 0, err
	}
	return total, os.WriteFile(RulesetPath, data, 0644)
}

func LoadRuleset() (*Ruleset, error) {
	data, err := os.ReadFile(RulesetPath)
	if err != nil {
		return nil, err
	}
	var rs Ruleset
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, err
	}
	return &rs, nil
}
