package sub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

func hysteria(password string, host string, inbound database.Inbound) string {
	port := strconv.Itoa(inbound.Port)

	params := url.Values{}

	if inbound.ObfsType != "" {
		params.Set("obfs", inbound.ObfsType)
		params.Set("obfs-password", inbound.ObfsPassword)
	}
	if inbound.HopEnabled && inbound.HopPort != "" {
		params.Set("mport", inbound.HopPort)
	}

	if inbound.ServerName != "" {
		params.Set("sni", inbound.ServerName)
	}
	if inbound.ALPN != "" {
		params.Set("alpn", inbound.ALPN)
	}
	if inbound.Insecure {
		params.Set("insecure", "1")
	}
	if inbound.ECHEnabled && inbound.ECHConfig != "" {
		params.Set("ech", extractECHConfig(inbound.ECHConfig))
	}

	name := url.PathEscape(inbound.Name)

	return "hysteria2://" + password + "@" + host + ":" + port + "?" + params.Encode() + "#" + name
}

func hysteriaClash(password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  - name: %s\n", inbound.Name)
	fmt.Fprintf(&sb, "    type: hysteria2\n")
	fmt.Fprintf(&sb, "    server: %s\n", host)
	fmt.Fprintf(&sb, "    port: %d\n", inbound.Port)
	if inbound.HopEnabled && inbound.HopPort != "" {
		fmt.Fprintf(&sb, "    ports: %s\n", inbound.HopPort)
		if inbound.HopInterval != "" {
			fmt.Fprintf(&sb, "    hop-interval: %s\n", inbound.HopInterval)
		}
	}
	fmt.Fprintf(&sb, "    password: %s\n", password)

	if inbound.ObfsType != "" {
		fmt.Fprintf(&sb, "    obfs: %s\n", inbound.ObfsType)
		fmt.Fprintf(&sb, "    obfs-password: %s\n", inbound.ObfsPassword)
	}

	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, "    sni: %s\n", inbound.ServerName)
	}
	if inbound.ALPN != "" {
		fmt.Fprintf(&sb, "    alpn:\n      - %s\n", inbound.ALPN)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, "    skip-cert-verify: true\n")
	}

	return sb.String()
}

func hysteriaSurge(password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s = hysteria2, %s, %d, password=%s", inbound.Name, host, inbound.Port, password)

	if inbound.ObfsType == "salamander" {
		fmt.Fprintf(&sb, ", salamander-password=%s", inbound.ObfsPassword)
	}

	if inbound.HopEnabled && inbound.HopPort != "" {
		fmt.Fprintf(&sb, ", port-hopping=%s", inbound.HopPort)
		if inbound.HopInterval != "" {
			fmt.Fprintf(&sb, ", port-hopping-interval=%s", inbound.HopInterval)
		}
	}

	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, ", sni=%s", inbound.ServerName)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, ", skip-cert-verify=true")
	}

	sb.WriteString("\n")
	return sb.String()
}

func hysteriaSingBox(password string, host string, inbound database.Inbound) string {
	out := map[string]any{
		"type":        "hysteria2",
		"tag":         "proxy",
		"server":      host,
		"server_port": inbound.Port,
		"password":    password,
	}

	if inbound.ObfsType != "" {
		out["obfs"] = map[string]any{
			"type":     inbound.ObfsType,
			"password": inbound.ObfsPassword,
		}
	}
	if inbound.HopEnabled && inbound.HopPort != "" {
		out["server_ports"] = []string{inbound.HopPort}
		parts := strings.SplitN(inbound.HopInterval, "-", 2)
		if len(parts) == 2 {
			out["hop_interval"] = strings.TrimSpace(parts[0]) + "s"
			out["hop_interval_max"] = strings.TrimSpace(parts[1]) + "s"
		}
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
	if inbound.ECHEnabled && inbound.ECHConfig != "" {
		tls["ech"] = map[string]any{
			"enabled": true,
			"config":  strings.Split(inbound.ECHConfig, "\n"),
		}
	}
	out["tls"] = tls

	data, _ := json.MarshalIndent(out, "        ", "    ")
	return string(data)
}
