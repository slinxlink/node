package template

import (
	"os"
	"strings"

	"github.com/slinxlink/node/internal/util"
)

const singBoxTemplate = `{
    "dns": {
        "servers": [
            {
                "tag": "remote",
                "address": "https://1.1.1.1/dns-query",
                "detour": "proxy"
            },
            {
                "tag": "local",
                "address": "https://223.5.5.5/dns-query",
                "detour": "direct"
            }
        ],
        "rules": [
            {
                "outbound": "any",
                "server": "local"
            },
            {
                "geosite": ["cn"],
                "server": "local"
            }
        ],
        "final": "remote"
    },
    "inbounds": [
        {
            "type": "tun",
            "tag": "tun-in",
            "address": ["172.19.0.1/30", "fdfe:dcba:9876::1/126"],
            "auto_route": true,
            "strict_route": true
        }
    ],
    "outbounds": [
        // __OUTBOUNDS__
        {
            "type": "direct",
            "tag": "direct"
        },
        {
            "type": "block",
            "tag": "block"
        },
        {
            "type": "dns",
            "tag": "dns-out"
        }
    ],
    "route": {
        "rules": [
            {
                "protocol": "dns",
                "outbound": "dns-out"
            },
            {
                "geosite": ["cn"],
                "outbound": "direct"
            },
            {
                "geoip": ["cn", "private"],
                "outbound": "direct"
            }
        ],
        "final": "proxy"
    }
}`

func GenerateSingBox() {
	if err := os.MkdirAll("data/ruleset", 0755); err != nil {
		util.Error("[singbox] 创建目录失败: %v", err)
		return
	}
	if err := os.WriteFile(SingBoxTemplatePath, []byte(singBoxTemplate), 0644); err != nil {
		util.Error("[singbox] 模板写入失败: %v", err)
		return
	}
	util.Info("[singbox] 模板已生成")
}

func RenderSingBox(outbounds []string) string {
	data, err := os.ReadFile(SingBoxTemplatePath)
	if err != nil {
		util.Error("[singbox] 读取模板失败: %v", err)
		return ""
	}

	outboundBlock := strings.Join(outbounds, ",\n    ")
	conf := strings.Replace(string(data), "// __OUTBOUNDS__", outboundBlock+",", 1)
	return conf
}
