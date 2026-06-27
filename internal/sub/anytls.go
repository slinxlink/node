package sub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

func anytls(password string, host string, inbound database.Inbound) string {
	port := strconv.Itoa(inbound.Port)
	params := url.Values{}

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
	params.Set("keepalive", fmt.Sprintf("%d,%d,%d",
		inbound.AnyTLSIdleSessionCheckInterval,
		inbound.AnyTLSIdleSessionTimeout,
		inbound.AnyTLSMinIdleSession,
	))

	name := url.PathEscape(inbound.Name)
	return "anytls://" + password + "@" + host + ":" + port + "?" + params.Encode() + "#" + name
}

func anytlsClash(password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  - name: %s\n", inbound.Name)
	fmt.Fprintf(&sb, "    type: anytls\n")
	fmt.Fprintf(&sb, "    server: %s\n", host)
	fmt.Fprintf(&sb, "    port: %d\n", inbound.Port)
	fmt.Fprintf(&sb, "    password: %s\n", password)
	fmt.Fprintf(&sb, "    udp: true\n")
	fmt.Fprintf(&sb, "    idle-session-check-interval: %d\n", inbound.AnyTLSIdleSessionCheckInterval)
	fmt.Fprintf(&sb, "    idle-session-timeout: %d\n", inbound.AnyTLSIdleSessionTimeout)
	fmt.Fprintf(&sb, "    min-idle-session: %d\n", inbound.AnyTLSMinIdleSession)

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

func anytlsSurge(password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s = anytls, %s, %d, password=%s", inbound.Name, host, inbound.Port, password)

	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, ", sni=%s", inbound.ServerName)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, ", skip-cert-verify=true")
	}
	if inbound.UTLS != "" {
		fmt.Fprintf(&sb, ", client-fingerprint=%s", inbound.UTLS)
	}

	sb.WriteString("\n")
	return sb.String()
}

func anytlsSingBox(password string, host string, inbound database.Inbound) string {
	out := map[string]any{
		"type":        "anytls",
		"tag":         "proxy",
		"server":      host,
		"server_port": inbound.Port,
		"password":    password,
	}

	if inbound.AnyTLSIdleSessionCheckInterval > 0 {
		out["idle_session_check_interval"] = fmt.Sprintf("%ds", inbound.AnyTLSIdleSessionCheckInterval)
	}
	if inbound.AnyTLSIdleSessionTimeout > 0 {
		out["idle_session_timeout"] = fmt.Sprintf("%ds", inbound.AnyTLSIdleSessionTimeout)
	}
	if inbound.AnyTLSMinIdleSession > 0 {
		out["min_idle_session"] = inbound.AnyTLSMinIdleSession
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
