package sub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

func tuic(uuid string, password string, host string, inbound database.Inbound) string {
	port := strconv.Itoa(inbound.Port)

	params := url.Values{}
	params.Set("congestion_control", inbound.TuicCongestionControl)

	if inbound.TuicUDPRelayMode != "" {
		params.Set("udp_relay_mode", inbound.TuicUDPRelayMode)
	}
	if inbound.ServerName != "" {
		params.Set("sni", inbound.ServerName)
	}
	if inbound.ALPN != "" {
		params.Set("alpn", inbound.ALPN)
	}
	if inbound.Insecure {
		params.Set("allow_insecure", "1")
	}

	name := url.PathEscape(inbound.Name)

	return "tuic://" + uuid + ":" + password + "@" + host + ":" + port + "?" + params.Encode() + "#" + name
}

func tuicClash(uuid string, password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  - name: %s\n", inbound.Name)
	fmt.Fprintf(&sb, "    type: tuic\n")
	fmt.Fprintf(&sb, "    server: %s\n", host)
	fmt.Fprintf(&sb, "    port: %d\n", inbound.Port)
	fmt.Fprintf(&sb, "    uuid: %s\n", uuid)
	fmt.Fprintf(&sb, "    password: %s\n", password)
	fmt.Fprintf(&sb, "    congestion-controller: %s\n", inbound.TuicCongestionControl)
	if inbound.TuicZeroRTT {
		fmt.Fprintf(&sb, "    reduce-rtt: true\n")
	}
	if inbound.TuicHeartbeat > 0 {
		fmt.Fprintf(&sb, "    heartbeat-interval: %d\n", inbound.TuicHeartbeat*1000)
	}
	if inbound.TuicUDPRelayMode != "" {
		fmt.Fprintf(&sb, "    udp-relay-mode: %s\n", inbound.TuicUDPRelayMode)
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

func tuicSurge(uuid string, password string, host string, inbound database.Inbound) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s = tuic-v5, %s, %d, uuid=%s, password=%s", inbound.Name, host, inbound.Port, uuid, password)

	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, ", sni=%s", inbound.ServerName)
	}
	if inbound.ALPN != "" {
		fmt.Fprintf(&sb, ", alpn=%s", inbound.ALPN)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, ", skip-cert-verify=true")
	}

	sb.WriteString("\n")
	return sb.String()
}

func tuicSingBox(uuid string, password string, host string, inbound database.Inbound) string {
	out := map[string]any{
		"type":               "tuic",
		"tag":                "proxy",
		"server":             host,
		"server_port":        inbound.Port,
		"uuid":               uuid,
		"password":           password,
		"congestion_control": inbound.TuicCongestionControl,
		"zero_rtt_handshake": inbound.TuicZeroRTT,
	}

	if inbound.TuicHeartbeat > 0 {
		out["heartbeat"] = fmt.Sprintf("%ds", inbound.TuicHeartbeat)
	}

	udpRelayMode := inbound.TuicUDPRelayMode
	if udpRelayMode == "" {
		udpRelayMode = "native"
	}
	out["udp_relay_mode"] = udpRelayMode

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
