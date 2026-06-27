package sub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

func vless(uuid string, host string, inbound database.Inbound) string {
	port := strconv.Itoa(inbound.Port)

	params := url.Values{}
	params.Set("encryption", "none")

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

	switch inbound.TLSType {
	case "tls":
		params.Set("security", "tls")
		if inbound.ServerName != "" {
			params.Set("sni", inbound.ServerName)
		}
		if inbound.UTLS != "" {
			params.Set("fp", inbound.UTLS)
		}
		if inbound.ALPN != "" {
			params.Set("alpn", inbound.ALPN)
		}
		if inbound.Insecure {
			params.Set("allowInsecure", "1")
		}
		if inbound.ECHEnabled && inbound.ECHConfig != "" {
			params.Set("ech", extractECHConfig(inbound.ECHConfig))
		}
	case "reality":
		params.Set("security", "reality")
		params.Set("sni", inbound.RealityServerName)
		var shortIDs []string
		if err := json.Unmarshal([]byte(inbound.RealityShortIDs), &shortIDs); err == nil && len(shortIDs) > 0 {
			params.Set("sid", shortIDs[0])
		}
		params.Set("pbk", inbound.RealityPublicKey)
		if inbound.UTLS != "" {
			params.Set("fp", inbound.UTLS)
		}
		if inbound.ALPN != "" {
			params.Set("alpn", inbound.ALPN)
		}
		if inbound.Flow != "" {
			params.Set("flow", inbound.Flow)
		}
	default:
		params.Set("security", "none")
	}

	name := url.PathEscape(inbound.Name)

	return "vless://" + uuid + "@" + host + ":" + port + "?" + params.Encode() + "#" + name
}

func vlessClash(uuid string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  - name: %s\n", inbound.Name)
	fmt.Fprintf(&sb, "    type: vless\n")
	fmt.Fprintf(&sb, "    server: %s\n", host)
	fmt.Fprintf(&sb, "    port: %d\n", inbound.Port)
	fmt.Fprintf(&sb, "    uuid: %s\n", uuid)

	// transport
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

	// tls
	switch inbound.TLSType {
	case "tls":
		fmt.Fprintf(&sb, "    tls: true\n")
		if inbound.ServerName != "" {
			fmt.Fprintf(&sb, "    servername: %s\n", inbound.ServerName)
		}
		if inbound.UTLS != "" {
			fmt.Fprintf(&sb, "    client-fingerprint: %s\n", inbound.UTLS)
		}
		if inbound.Insecure {
			fmt.Fprintf(&sb, "    skip-cert-verify: true\n")
		}
	case "reality":
		fmt.Fprintf(&sb, "    tls: true\n")
		fmt.Fprintf(&sb, "    servername: %s\n", inbound.RealityServerName)
		if inbound.UTLS != "" {
			fmt.Fprintf(&sb, "    client-fingerprint: %s\n", inbound.UTLS)
		}
		if inbound.Flow != "" {
			fmt.Fprintf(&sb, "    flow: %s\n", inbound.Flow)
		}
		var shortIDs []string
		if err := json.Unmarshal([]byte(inbound.RealityShortIDs), &shortIDs); err == nil && len(shortIDs) > 0 {
			fmt.Fprintf(&sb, "    reality-opts:\n")
			fmt.Fprintf(&sb, "      public-key: %s\n", inbound.RealityPublicKey)
			fmt.Fprintf(&sb, "      short-id: %s\n", shortIDs[0])
		}
	}

	return sb.String()
}

func vlessSingBox(uuid string, host string, inbound database.Inbound) string {
	out := map[string]any{
		"type":        "vless",
		"tag":         "proxy",
		"server":      host,
		"server_port": inbound.Port,
		"uuid":        uuid,
		"network":     "tcp",
	}

	if inbound.Flow != "" {
		out["flow"] = inbound.Flow
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

	switch inbound.TLSType {
	case "tls":
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
		if inbound.ECHEnabled && inbound.ECHConfig != "" {
			tls["ech"] = map[string]any{
				"enabled": true,
				"config":  strings.Split(inbound.ECHConfig, "\n"),
			}
		}
		out["tls"] = tls
	case "reality":
		var shortIDs []string
		json.Unmarshal([]byte(inbound.RealityShortIDs), &shortIDs)
		shortID := ""
		if len(shortIDs) > 0 {
			shortID = shortIDs[0]
		}
		out["tls"] = map[string]any{
			"enabled":     true,
			"server_name": inbound.RealityServerName,
			"utls": map[string]any{
				"enabled":     true,
				"fingerprint": inbound.UTLS,
			},
			"reality": map[string]any{
				"enabled":    true,
				"public_key": inbound.RealityPublicKey,
				"short_id":   shortID,
			},
		}
	}

	data, _ := json.MarshalIndent(out, "        ", "    ")
	return string(data)
}
