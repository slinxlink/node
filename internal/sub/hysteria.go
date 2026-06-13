package sub

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

func hysteria(password string, host string, inbound database.Inbound) string {
	port := strconv.Itoa(inbound.Port)

	params := url.Values{}

	if inbound.ServerName != "" {
		params.Set("sni", inbound.ServerName)
	}
	if inbound.ALPN != "" {
		params.Set("alpn", inbound.ALPN)
	}
	if inbound.UTLS != "" {
		params.Set("fp", inbound.UTLS)
	}
	if inbound.Insecure {
		params.Set("insecure", "1")
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
	fmt.Fprintf(&sb, "    password: %s\n", password)

	if inbound.ServerName != "" {
		fmt.Fprintf(&sb, "    sni: %s\n", inbound.ServerName)
	}
	if inbound.ALPN != "" {
		fmt.Fprintf(&sb, "    alpn:\n      - %s\n", inbound.ALPN)
	}
	if inbound.UTLS != "" {
		fmt.Fprintf(&sb, "    client-fingerprint: %s\n", inbound.UTLS)
	}
	if inbound.Insecure {
		fmt.Fprintf(&sb, "    skip-cert-verify: true\n")
	}

	return sb.String()
}
