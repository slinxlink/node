package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/slinxlink/node/internal/util"
)

func GenerateSurge() {
	rs, err := LoadRuleset()
	if err != nil {
		util.Error("[surge] 读取规则集失败: %v", err)
		return
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `# SLINX Node · Surge
[General]
loglevel = notify
dns-server = 8.8.8.8, 1.1.1.1
skip-proxy = 127.0.0.1, 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, 100.64.0.0/10, localhost, *.local

[Proxy]
# __PROXIES__

[Proxy Group]
PROXY = select, # __PROXY_NAMES__
`)

	for _, g := range rs.Groups {
		if g.Name == "_final" {
			continue
		}
		defaultFirst := g.Name == "APPLE" || g.Name == "MICROSOFT"
		if defaultFirst {
			fmt.Fprintf(&sb, "%s = select, DIRECT, PROXY, # __PROXY_NAMES__\n", g.Name)
		} else {
			fmt.Fprintf(&sb, "%s = select, PROXY, DIRECT, # __PROXY_NAMES__\n", g.Name)
		}
	}

	fmt.Fprintf(&sb, `DEFAULT = select, PROXY, DIRECT, # __PROXY_NAMES__

[Rule]
`)

	for _, g := range rs.Groups {
		if g.Name == "_final" {
			fmt.Fprintf(&sb, "GEOIP,CN,DIRECT\n")
			fmt.Fprintf(&sb, "FINAL,DEFAULT\n")
			continue
		}
		for _, rule := range g.Rules {
			if strings.HasPrefix(rule, "IP-CIDR") {
				fmt.Fprintf(&sb, "%s,%s,no-resolve\n", rule, g.Name)
			} else {
				fmt.Fprintf(&sb, "%s,%s\n", rule, g.Name)
			}
		}
	}

	if err := os.WriteFile(SurgeTemplatePath, []byte(sb.String()), 0644); err != nil {
		util.Error("[surge] 模板写入失败: %v", err)
		return
	}
	util.Info("[surge] 模板已更新")
}

func RenderSurge(proxies []string, proxyNames []string) string {
	data, err := os.ReadFile(SurgeTemplatePath)
	if err != nil {
		util.Error("[surge] 读取模板失败: %v", err)
		return ""
	}

	proxyBlock := strings.Join(proxies, "\n")
	nameBlock := strings.Join(proxyNames, ", ")

	conf := string(data)
	conf = strings.Replace(conf, "# __PROXIES__", proxyBlock, 1)
	conf = strings.ReplaceAll(conf, "# __PROXY_NAMES__", nameBlock)

	return conf
}
