package sub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

type vmessLink struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	ALPN string `json:"alpn"`
	FP   string `json:"fp"`
}

func vmess(uuid string, host string, inbound database.Inbound) string {

	link := vmessLink{
		V:    "2",
		Ps:   inbound.Name,
		Add:  host,
		Port: strconv.Itoa(inbound.Port),
		ID:   uuid,
		Aid:  "0",
		Scy:  "auto",
		Type: "none",
	}

	switch inbound.Transport {
	case "websocket":
		link.Net = "ws"
		link.Path = inbound.WsPath
		link.Host = inbound.WsHost
	default:
		link.Net = "tcp"
	}

	switch inbound.TLSType {
	case "tls":
		link.TLS = "tls"
		link.SNI = inbound.ServerName
		link.ALPN = inbound.ALPN
		link.FP = inbound.UTLS
	default:
		link.TLS = "none"
	}

	data, _ := json.Marshal(link)
	return "vmess://" + base64.StdEncoding.EncodeToString(data)
}

func vmessClash(uuid string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  - name: %s\n", inbound.Name)
	fmt.Fprintf(&sb, "    type: vmess\n")
	fmt.Fprintf(&sb, "    server: %s\n", host)
	fmt.Fprintf(&sb, "    port: %d\n", inbound.Port)
	fmt.Fprintf(&sb, "    uuid: %s\n", uuid)
	fmt.Fprintf(&sb, "    alterId: 0\n")
	fmt.Fprintf(&sb, "    cipher: auto\n")
	fmt.Fprintf(&sb, "    udp: true\n")

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

	switch inbound.TLSType {
	case "tls":
		fmt.Fprintf(&sb, "    tls: true\n")
		if inbound.ServerName != "" {
			fmt.Fprintf(&sb, "    servername: %s\n", inbound.ServerName)
		}
		if inbound.ALPN != "" {
			fmt.Fprintf(&sb, "    alpn:\n      - %s\n", inbound.ALPN)
		}
		if inbound.UTLS != "" {
			fmt.Fprintf(&sb, "    client-fingerprint: %s\n", inbound.UTLS)
		}
	}

	return sb.String()
}

func vmessSurge(uuid string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s = vmess, %s, %d, username=%s, vmess-aead=true", inbound.Name, host, inbound.Port, uuid)

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

	switch inbound.TLSType {
	case "tls":
		fmt.Fprintf(&sb, ", tls=true")
		if inbound.ServerName != "" {
			fmt.Fprintf(&sb, ", sni=%s", inbound.ServerName)
		}
		if inbound.UTLS != "" {
			fmt.Fprintf(&sb, ", client-fingerprint=%s", inbound.UTLS)
		}
		if inbound.Insecure {
			fmt.Fprintf(&sb, ", skip-cert-verify=true")
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func vmessSingBox(uuid string, host string, inbound database.Inbound) string {
	out := map[string]any{
		"type":        "vmess",
		"tag":         "proxy",
		"server":      host,
		"server_port": inbound.Port,
		"uuid":        uuid,
		"security":    "auto",
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

	if inbound.TLSType == "tls" {
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
	}

	data, _ := json.MarshalIndent(out, "        ", "    ")
	return string(data)
}
