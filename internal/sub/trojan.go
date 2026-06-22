package sub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

func trojan(password string, host string, inbound database.Inbound) string {
	port := strconv.Itoa(inbound.Port)
	params := url.Values{}

	switch inbound.Transport {
	case "websocket":
		params.Set("type", "ws")
		if inbound.WsPath != "" {
			params.Set("path", inbound.WsPath)
		}
		if inbound.WsHost != "" {
			params.Set("host", inbound.WsHost)
		}
	default:
		params.Set("type", "tcp")
	}

	params.Set("security", "tls")
	if inbound.ServerName != "" {
		params.Set("sni", inbound.ServerName)
	}
	if inbound.UTLS != "" {
		params.Set("fp", inbound.UTLS)
	}
	if inbound.Insecure {
		params.Set("allowInsecure", "1")
	}
	if inbound.ALPN != "" {
		params.Set("alpn", inbound.ALPN)
	}

	name := url.PathEscape(inbound.Name)
	return "trojan://" + password + "@" + host + ":" + port + "?" + params.Encode() + "#" + name
}

func trojanClash(password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  - name: %s\n", inbound.Name)
	fmt.Fprintf(&sb, "    type: trojan\n")
	fmt.Fprintf(&sb, "    server: %s\n", host)
	fmt.Fprintf(&sb, "    port: %d\n", inbound.Port)
	fmt.Fprintf(&sb, "    password: %s\n", password)

	switch inbound.Transport {
	case "websocket":
		fmt.Fprintf(&sb, "    network: ws\n")
		if inbound.WsPath != "" || inbound.WsHost != "" {
			fmt.Fprintf(&sb, "    ws-opts:\n")
			if inbound.WsPath != "" {
				fmt.Fprintf(&sb, "      path: %s\n", inbound.WsPath)
			}
			if inbound.WsHost != "" {
				fmt.Fprintf(&sb, "      headers:\n")
				fmt.Fprintf(&sb, "        Host: %s\n", inbound.WsHost)
			}
		}
	default:
		fmt.Fprintf(&sb, "    network: tcp\n")
	}

	fmt.Fprintf(&sb, "    tls: true\n")
	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, "    sni: %s\n", inbound.ServerName)
	}
	if inbound.UTLS != "" {
		fmt.Fprintf(&sb, "    client-fingerprint: %s\n", inbound.UTLS)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, "    skip-cert-verify: true\n")
	}
	if inbound.ALPN != "" {
		fmt.Fprintf(&sb, "    alpn:\n")
		for _, a := range strings.Split(inbound.ALPN, ",") {
			fmt.Fprintf(&sb, "      - %s\n", strings.TrimSpace(a))
		}
	}

	return sb.String()
}

func trojanSurge(password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s = trojan, %s, %d, password=%s", inbound.Name, host, inbound.Port, password)
	fmt.Fprintf(&sb, ", tls=true")

	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, ", sni=%s", inbound.ServerName)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, ", skip-cert-verify=true")
	}
	if inbound.UTLS != "" {
		fmt.Fprintf(&sb, ", client-fingerprint=%s", inbound.UTLS)
	}

	switch inbound.Transport {
	case "websocket":
		fmt.Fprintf(&sb, ", ws=true")
		if inbound.WsPath != "" {
			fmt.Fprintf(&sb, ", ws-path=%s", inbound.WsPath)
		}
		if inbound.WsHost != "" {
			fmt.Fprintf(&sb, ", ws-headers=Host:%s", inbound.WsHost)
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func trojanSingBox(password string, host string, inbound database.Inbound) string {
	out := map[string]any{
		"type":        "trojan",
		"tag":         "proxy",
		"server":      host,
		"server_port": inbound.Port,
		"password":    password,
		"network":     "tcp",
	}

	if inbound.Transport == "websocket" {
		t := map[string]any{"type": "ws"}
		if inbound.WsPath != "" {
			t["path"] = inbound.WsPath
		}
		if inbound.WsHost != "" {
			t["headers"] = map[string]string{"Host": inbound.WsHost}
		}
		out["transport"] = t
	}

	tls := map[string]any{"enabled": true}
	if inbound.ServerName != "" {
		tls["server_name"] = inbound.ServerName
	}
	if inbound.UTLS != "" {
		tls["utls"] = map[string]any{"enabled": true, "fingerprint": inbound.UTLS}
	}
	if inbound.Insecure {
		tls["insecure"] = true
	}
	if inbound.ALPN != "" {
		tls["alpn"] = strings.Split(inbound.ALPN, ",")
	}
	if inbound.CipherSuites != "" {
		tls["cipher_suites"] = strings.Split(inbound.CipherSuites, ",")
	}
	if inbound.TLSMinVersion != "" {
		tls["min_version"] = inbound.TLSMinVersion
	}
	if inbound.TLSMaxVersion != "" {
		tls["max_version"] = inbound.TLSMaxVersion
	}
	out["tls"] = tls

	data, _ := json.MarshalIndent(out, "        ", "    ")
	return string(data)
}
