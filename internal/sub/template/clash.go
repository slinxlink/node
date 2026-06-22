package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/slinxlink/node/internal/util"
)

var clashGroupName = map[string]string{
	"AI":        "🤖 AI",
	"MEDIA":     "🎬 媒体",
	"SOCIAL":    "💬 通讯",
	"GAME":      "🎮 游戏",
	"APPLE":     "🍎 苹果",
	"GOOGLE":    "🔍 谷歌",
	"MICROSOFT": "Ⓜ️ 微软",
}

func GenerateClash() {
	rs, err := LoadRuleset()
	if err != nil {
		util.Error("[clash] 读取规则集失败: %v", err)
		return
	}

	var sb strings.Builder
	sb.WriteString(`# SLINX Node · Clash 订阅
mixed-port: 7890
allow-lan: true
mode: rule
log-level: error
external-controller: 0.0.0.0:9090

proxies: []

proxy-groups:
  - name: 🚀 代理
    type: select
    proxies:
# __PROXIES__
`)

	for _, g := range rs.Groups {
		if g.Name == "_final" {
			continue
		}
		name := g.Name
		if cn, ok := clashGroupName[g.Name]; ok {
			name = cn
		}
		defaultFirst := g.Name == "APPLE" || g.Name == "MICROSOFT"
		fmt.Fprintf(&sb, "  - name: %s\n", name)
		fmt.Fprintf(&sb, "    type: select\n")
		fmt.Fprintf(&sb, "    proxies:\n")
		if defaultFirst {
			fmt.Fprintf(&sb, "      - 🎯 直连\n      - 🚀 代理\n")
		} else {
			fmt.Fprintf(&sb, "      - 🚀 代理\n      - 🎯 直连\n")
		}
		fmt.Fprintf(&sb, "# __PROXIES__\n")
	}

	fmt.Fprintf(&sb, `  - name: 🐟 默认
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

	for _, g := range rs.Groups {
		if g.Name == "_final" {
			fmt.Fprintf(&sb, "  - GEOIP,LAN,🎯 直连,no-resolve\n")
			fmt.Fprintf(&sb, "  - GEOSITE,CN,🎯 直连\n")
			fmt.Fprintf(&sb, "  - GEOIP,CN,🎯 直连,no-resolve\n")
			fmt.Fprintf(&sb, "  - MATCH,🐟 默认\n")
			continue
		}
		name := g.Name
		if cn, ok := clashGroupName[g.Name]; ok {
			name = cn
		}
		for _, rule := range g.Rules {
			if strings.HasPrefix(rule, "IP-CIDR") {
				fmt.Fprintf(&sb, "  - %s,%s,no-resolve\n", rule, name)
			} else {
				fmt.Fprintf(&sb, "  - %s,%s\n", rule, name)
			}
		}
	}

	if err := os.WriteFile(ClashTemplatePath, []byte(sb.String()), 0644); err != nil {
		util.Error("[clash] 模板写入失败: %v", err)
		return
	}
	util.Info("[clash] 模板已更新")
}

func RenderClash(proxies []string) string {
	data, err := os.ReadFile(ClashTemplatePath)
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
