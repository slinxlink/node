package core

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

// ── sing-box config structs ──────────────────────────────────────────────────

type config struct {
	Log          log           `json:"log"`
	Experimental *experimental `json:"experimental,omitempty"`
	Inbounds     []inbounds    `json:"inbounds"`
	Outbounds    []outbounds   `json:"outbounds"`
	Endpoints    []endpoints   `json:"endpoints,omitempty"`
	Route        route         `json:"route"`
}

type log struct {
	Level     string `json:"level"`
	Output    string `json:"output,omitempty"`
	Timestamp bool   `json:"timestamp,omitempty"`
}

type experimental struct {
	V2RayAPI *v2rayAPI `json:"v2ray_api,omitempty"`
	ClashAPI *clashAPI `json:"clash_api,omitempty"`
}

type v2rayAPI struct {
	Listen string `json:"listen"`
	Stats  stats  `json:"stats"`
}

type stats struct {
	Enabled   bool     `json:"enabled"`
	Inbounds  []string `json:"inbounds,omitempty"`
	Outbounds []string `json:"outbounds,omitempty"`
	Users     []string `json:"users,omitempty"`
}

type clashAPI struct {
	ExternalController string `json:"external_controller"`
	Secret             string `json:"secret,omitempty"`
}

type inbounds struct {
	Type       string      `json:"type"`
	Tag        string      `json:"tag"`
	Listen     string      `json:"listen"`
	Port       int         `json:"listen_port"`
	UDPTimeout string      `json:"udp_timeout,omitempty"`
	Transport  *transport  `json:"transport,omitempty"`
	Masquerade interface{} `json:"masquerade,omitempty"`
	Obfs       *obfs       `json:"obfs,omitempty"`
	TLS        *tls        `json:"tls,omitempty"`
	Users      []user      `json:"users,omitempty"`
}

type transport struct {
	Type    string            `json:"type"`
	Path    string            `json:"path,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type obfs struct {
	Type          string `json:"type"`
	Password      string `json:"password"`
	MinPacketSize int    `json:"min_packet_size,omitempty"`
	MaxPacketSize int    `json:"max_packet_size,omitempty"`
}

type tls struct {
	Enabled         bool     `json:"enabled"`
	ServerName      string   `json:"server_name,omitempty"`
	ALPN            []string `json:"alpn,omitempty"`
	MinVersion      string   `json:"min_version,omitempty"`
	MaxVersion      string   `json:"max_version,omitempty"`
	CipherSuites    []string `json:"cipher_suites,omitempty"`
	CertificatePath string   `json:"certificate_path,omitempty"`
	KeyPath         string   `json:"key_path,omitempty"`
	Reality         *reality `json:"reality,omitempty"`
}

type reality struct {
	Enabled           bool      `json:"enabled"`
	Handshake         handshake `json:"handshake"`
	PrivateKey        string    `json:"private_key"`
	ShortID           []string  `json:"short_id"`
	MaxTimeDifference string    `json:"max_time_difference,omitempty"`
}

type handshake struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
}

type user struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid,omitempty"`
	Password string `json:"password,omitempty"`
	Flow     string `json:"flow,omitempty"`
}

type outbounds struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

type endpoints struct {
	Type       string         `json:"type"`
	Tag        string         `json:"tag"`
	MTU        int            `json:"mtu,omitempty"`
	Address    []string       `json:"address,omitempty"`
	PrivateKey string         `json:"private_key,omitempty"`
	Peers      []endpointPeer `json:"peers,omitempty"`
	UDPTimeout string         `json:"udp_timeout,omitempty"`
	Workers    int            `json:"workers,omitempty"`
}

type endpointPeer struct {
	Address    string   `json:"address"`
	Port       int      `json:"port"`
	PublicKey  string   `json:"public_key"`
	AllowedIPs []string `json:"allowed_ips,omitempty"`
	Reserved   []int    `json:"reserved,omitempty"`
}

type route struct {
	Rules []map[string]any `json:"rules"`
	Final string           `json:"final,omitempty"`
}

type routeRule struct {
	IPIsPrivate *bool  `json:"ip_is_private,omitempty"`
	Outbound    string `json:"outbound,omitempty"`
}

// ── database loader ──────────────────────────────────────────────────────────

type db struct {
	Core       database.Core
	Inbounds   []database.Inbound
	Users      []database.User
	Boards     []database.Board
	BoardUsers map[uint][]database.BoardUser
	Endpoints  []database.Endpoint
	Rules      []database.Rule
	UserNames  []string
}

func loadDatabase() (db, error) {
	var d db
	database.DB.First(&d.Core)
	database.DB.Where("enable = ?", true).Find(&d.Inbounds)
	database.DB.Where("enable = ?", true).Find(&d.Users)
	database.DB.Where("enable = ?", true).Find(&d.Endpoints)
	database.DB.Order("sort asc, `index` asc").Find(&d.Rules)

	var cfg database.Config
	database.DB.First(&cfg)

	if cfg.BoardEnable {
		database.DB.Where("enable = ?", true).Find(&d.Boards)
		d.BoardUsers = make(map[uint][]database.BoardUser)
		for _, b := range d.Boards {
			var bu []database.BoardUser
			database.DB.Where("board_id = ?", b.ID).Find(&bu)
			d.BoardUsers[b.ID] = bu
		}
	}

	for _, u := range d.Users {
		d.UserNames = append(d.UserNames, u.Name)
	}
	for _, bus := range d.BoardUsers {
		for _, bu := range bus {
			d.UserNames = append(d.UserNames, fmt.Sprintf("BoardUser_%d", bu.UserID))
		}
	}

	return d, nil
}

// ── config generator ─────────────────────────────────────────────────────────

func generateConfig() error {
	d, err := loadDatabase()
	if err != nil {
		return err
	}

	cfg := config{
		Log: log{
			Level: d.Core.LogLevel,
			Output: func() string {
				if d.Core.LogEnable {
					return d.Core.LogPath
				}
				return ""
			}(),
			Timestamp: true,
		},
		Experimental: &experimental{
			V2RayAPI: &v2rayAPI{
				Listen: "127.0.0.1:2048",
				Stats: stats{
					Enabled: true,
					Users:   d.UserNames,
				},
			},
			ClashAPI: &clashAPI{
				ExternalController: "127.0.0.1:9090",
			},
		},
	}

	for _, ib := range d.Inbounds {
		var ibUsers []database.User
		for _, u := range d.Users {
			var ids []int
			json.Unmarshal([]byte(u.Inbounds), &ids)
			for _, id := range ids {
				if id == int(ib.ID) {
					ibUsers = append(ibUsers, u)
					break
				}
			}
		}

		var boardUsers []database.BoardUser
		for _, b := range d.Boards {
			if b.Inbound != int(ib.ID) {
				continue
			}
			boardUsers = append(boardUsers, d.BoardUsers[b.ID]...)
		}

		users := buildUsers(ib.Protocol, ibUsers, boardUsers, ib.Flow)
		ic, err := buildInbound(ib, users)
		if err != nil {
			continue
		}

		cfg.Inbounds = append(cfg.Inbounds, ic)
	}

	cfg.Outbounds = []outbounds{
		{
			Type: "direct",
			Tag:  "direct",
		},
		{
			Type: "block",
			Tag:  "block",
		},
	}

	for _, ep := range d.Endpoints {
		ec, err := buildEndpoint(ep)
		if err != nil {
			continue
		}
		cfg.Endpoints = append(cfg.Endpoints, ec)
	}

	cfg.Route = route{
		Rules: buildRouteRules(d.Rules),
		Final: "direct",
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(Default.ConfigPath, data, 0644)
}

// ── inbound builders ─────────────────────────────────────────────────────────

func buildInbound(ib database.Inbound, users []user) (inbounds, error) {
	switch ib.Protocol {
	case "hysteria":
		return buildHysteria(ib, users)
	case "vless":
		return buildVless(ib, users)
	case "vmess":
		return buildVmess(ib, users)
	}
	return inbounds{}, fmt.Errorf("unsupported protocol: %s", ib.Protocol)
}

func buildHysteria(ib database.Inbound, users []user) (inbounds, error) {
	ic := buildBase(ib)
	ic.UDPTimeout = fmt.Sprintf("%ds", ib.UDPTimeout)
	ic.Masquerade = buildMasquerade(ib)
	ic.Obfs = buildObfs(ib)
	ic.TLS = buildTLS(ib)
	ic.Users = users
	return ic, nil
}

func buildVless(ib database.Inbound, users []user) (inbounds, error) {
	ic := buildBase(ib)
	ic.Transport = buildTransport(ib)
	ic.TLS = buildTLS(ib)
	ic.Users = users
	return ic, nil
}

func buildVmess(ib database.Inbound, users []user) (inbounds, error) {
	ic := buildBase(ib)
	ic.Transport = buildTransport(ib)
	ic.TLS = buildTLS(ib)
	ic.Users = users
	return ic, nil
}

func buildBase(ib database.Inbound) inbounds {
	protocol := ib.Protocol
	if protocol == "hysteria" {
		protocol = "hysteria2"
	}
	return inbounds{
		Type:   protocol,
		Tag:    fmt.Sprintf("%d", ib.Port),
		Listen: "::",
		Port:   ib.Port,
	}
}

// ── field builders ───────────────────────────────────────────────────────────

func buildTransport(ib database.Inbound) *transport {
	switch ib.Transport {
	case "websocket":
		t := &transport{
			Type: "ws",
			Path: ib.WsPath,
		}
		if ib.WsHost != "" {
			t.Headers = map[string]string{"Host": ib.WsHost}
		}
		return t
	default:
		return nil
	}
}

func buildMasquerade(ib database.Inbound) interface{} {
	if !ib.MasqueradeEnabled {
		return nil
	}
	switch ib.MasqueradeType {
	case "proxy":
		return map[string]interface{}{
			"type":         "proxy",
			"url":          ib.MasqueradeURL,
			"rewrite_host": ib.RewriteHost,
		}
	case "file":
		return map[string]interface{}{
			"type":      "file",
			"directory": ib.MasqueradePath,
		}
	case "string":
		return map[string]interface{}{
			"type":        "string",
			"status_code": ib.MasqueradeCode,
			"content":     ib.MasqueradeBody,
		}
	default:
		return nil
	}
}

func buildObfs(ib database.Inbound) *obfs {
	if ib.ObfsType == "" {
		return nil
	}
	o := &obfs{
		Type:     ib.ObfsType,
		Password: ib.ObfsPassword,
	}
	if ib.ObfsType == "gecko" {
		o.MinPacketSize = ib.ObfsMinPacketSize
		o.MaxPacketSize = ib.ObfsMaxPacketSize
	}
	return o
}

func buildTLS(ib database.Inbound) *tls {
	switch ib.TLSType {
	case "tls":
		var ids []int
		json.Unmarshal([]byte(ib.Certs), &ids)
		if len(ids) == 0 || ids[0] == 0 {
			return nil
		}
		var cert database.Cert
		if database.DB.First(&cert, ids[0]).Error != nil {
			return nil
		}
		t := &tls{
			Enabled:         true,
			ServerName:      ib.ServerName,
			CertificatePath: cert.CertPath,
			KeyPath:         cert.KeyPath,
		}
		if ib.ALPN != "" {
			t.ALPN = strings.Split(ib.ALPN, ",")
		}
		if ib.TLSMinVersion != "" {
			t.MinVersion = ib.TLSMinVersion
		}
		if ib.TLSMaxVersion != "" {
			t.MaxVersion = ib.TLSMaxVersion
		}
		if ib.CipherSuites != "" {
			t.CipherSuites = strings.Split(ib.CipherSuites, ",")
		}
		return t
	case "reality":
		var shortIDs []string
		json.Unmarshal([]byte(ib.RealityShortIDs), &shortIDs)
		var maxTimeDiff string
		if ib.RealityMaxTimeDiff > 0 {
			maxTimeDiff = fmt.Sprintf("%dms", ib.RealityMaxTimeDiff)
		}
		return &tls{
			Enabled:    true,
			ServerName: ib.RealityServerName,
			Reality: &reality{
				Enabled: true,
				Handshake: handshake{
					Server:     ib.RealityServer,
					ServerPort: ib.RealityServerPort,
				},
				PrivateKey:        ib.RealityPrivateKey,
				ShortID:           shortIDs,
				MaxTimeDifference: maxTimeDiff,
			},
		}
	default:
		return nil
	}
}

// ── user builders ────────────────────────────────────────────────────────────

func buildUsers(protocol string, users []database.User, boardUsers []database.BoardUser, flow string) []user {
	var result []user
	for _, u := range users {
		if !u.Enable {
			continue
		}
		result = append(result, buildUser(protocol, u.Name, u.UUID, u.Password, flow))
	}
	for _, u := range boardUsers {
		result = append(result, buildUser(protocol, fmt.Sprintf("BoardUser_%d", u.UserID), u.UUID, u.Passwd, flow))
	}
	return result
}

func buildUser(protocol, name, uuid, password, flow string) user {
	u := user{Name: name}
	switch protocol {
	case "vless", "vmess":
		u.UUID = uuid
		u.Flow = flow
	case "hysteria":
		u.Password = password
	}
	return u
}

// ── endpoint builders ────────────────────────────────────────────────────────────

func buildEndpoint(ep database.Endpoint) (endpoints, error) {
	switch ep.Type {
	case "wireguard":
		return buildWireguardEndpoint(ep)
	}
	return endpoints{}, fmt.Errorf("unsupported endpoint type: %s", ep.Type)
}

func buildWireguardEndpoint(ep database.Endpoint) (endpoints, error) {
	var address []string
	json.Unmarshal([]byte(ep.Address), &address)

	var reserved []int
	json.Unmarshal([]byte(ep.Reserved), &reserved)

	var allowedIPs []string
	if ep.AllowedIPs != "" {
		allowedIPs = strings.Split(ep.AllowedIPs, ",")
	}

	var udpTimeout string
	if ep.UDPTimeout > 0 {
		udpTimeout = fmt.Sprintf("%ds", ep.UDPTimeout)
	}

	ec := endpoints{
		Type:       "wireguard",
		Tag:        ep.Tag,
		MTU:        ep.MTU,
		Address:    address,
		PrivateKey: ep.PrivateKey,
		Peers: []endpointPeer{
			{
				Address:    ep.PeerAddress,
				Port:       ep.PeerPort,
				PublicKey:  ep.PeerPublicKey,
				AllowedIPs: allowedIPs,
				Reserved:   reserved,
			},
		},
		UDPTimeout: udpTimeout,
		Workers:    ep.Workers,
	}
	return ec, nil
}

func buildRouteRules(rules []database.Rule) []map[string]any {
	result := []map[string]any{
		{"ip_is_private": true, "outbound": "block"},
	}

	grouped := map[int]map[string]any{}
	var sorts []int
	for _, r := range rules {
		if _, ok := grouped[r.Sort]; !ok {
			grouped[r.Sort] = map[string]any{}
			sorts = append(sorts, r.Sort)
		}
		var v any
		if err := json.Unmarshal([]byte(r.Value), &v); err == nil {
			grouped[r.Sort][r.Key] = v
		}
	}
	sort.Ints(sorts)

	for _, s := range sorts {
		result = append(result, grouped[s])
	}
	return result
}
