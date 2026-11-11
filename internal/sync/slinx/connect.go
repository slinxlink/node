package slinx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
	"gorm.io/gorm"
)

var ErrSave = errors.New("保存失败")

type nodeConfig struct {
	Protocol     string `json:"protocol"`
	Port         int    `json:"port"`
	Network      string `json:"network"`
	Host         string `json:"host"`
	Path         string `json:"path"`
	ObfsPassword string `json:"obfs_password"`
	Security     string `json:"security"`
	SNI          string `json:"sni"`
	Fingerprint  string `json:"fingerprint"`
	ShortID      string `json:"short_id"`
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`
}

const masqueradeURL = "https://www.bing.com"

func Connect(host string, nodeID int, key string) (database.Board, database.Inbound, error) {
	var board database.Board
	var ib database.Inbound

	// ── 1. 检查参数 ──────────────────────────────────────────────────────
	host = strings.TrimRight(strings.TrimSpace(host), "/")
	if host == "" {
		return board, ib, errors.New("面板地址不能为空")
	}
	if nodeID <= 0 {
		return board, ib, errors.New("节点 ID 不能为 0")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return board, ib, errors.New("通讯密钥不能为空")
	}

	var existing database.Board
	if database.DB.Where(database.Board{Host: host, NodeID: nodeID}).First(&existing).Error == nil {
		return board, ib, errors.New("这个节点已经对接过了")
	}

	// ── 2. 获取配置 ────────────────────────────────────────────────────────
	cfg, err := fetchNodeConfig(host, nodeID, key)
	if err != nil {
		return board, ib, err
	}

	// ── 3. 生成入站 ──────────────────────────────────────────────────
	ib, err = buildInbound(cfg, nodeID)
	if err != nil {
		return board, ib, err
	}

	// ── 4. 检查端口和证书 ────────────────────────────────────────────────
	if msg := util.ValidatePort(ib.Port, database.UsedPorts()); msg != "" {
		return board, ib, fmt.Errorf("端口 %d：%s", ib.Port, msg)
	}
	if ib.TLSType == "tls" {
		certID, err := findCert(ib.ServerName)
		if err != nil {
			return board, ib, err
		}
		ib.Certs = fmt.Sprintf("[%d]", certID)
	}

	// ── 5. 入站和对接一起写库 ────────────────────────────────────────────
	board = database.Board{
		Enable:       true,
		Name:         strings.TrimPrefix(strings.TrimPrefix(host, "https://"), "http://"),
		Host:         host,
		NodeID:       nodeID,
		Key:          key,
		Type:         "SLINX",
		SyncInterval: 60,
	}
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ib).Error; err != nil {
			return err
		}
		board.Inbound = int(ib.ID)
		return tx.Create(&board).Error
	})
	if err != nil {
		return board, ib, fmt.Errorf("%w: %v", ErrSave, err)
	}

	util.Info("[board] 一键对接: %s 节点 %d，入站 %s:%d", board.Name, board.NodeID, ib.Protocol, ib.Port)
	return board, ib, nil
}

func fetchNodeConfig(host string, nodeID int, key string) (nodeConfig, error) {
	var cfg nodeConfig

	req, err := http.NewRequest(http.MethodGet, host+"/api/node", nil)
	if err != nil {
		return cfg, fmt.Errorf("面板地址不正确: %w", err)
	}
	req.Header.Set("Node-Id", fmt.Sprintf("%d", nodeID))
	req.Header.Set("Node-Key", key)

	res, err := client.Do(req)
	if err != nil {
		return cfg, fmt.Errorf("连不上面板: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var e struct {
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&e)
		return cfg, fmt.Errorf("面板返回 %d %s", res.StatusCode, e.Message)
	}

	if err := json.NewDecoder(res.Body).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("解析面板配置失败: %w", err)
	}
	return cfg, nil
}

func buildInbound(cfg nodeConfig, nodeID int) (database.Inbound, error) {
	ib := database.Inbound{
		Enable:   true,
		Name:     fmt.Sprintf("SLINX #%d", nodeID),
		Protocol: cfg.Protocol,
		Port:     cfg.Port,

		Transport: "raw",
		WsPath:    "/",

		UDPTimeout:        60,
		MasqueradeType:    "default",
		MasqueradeCode:    200,
		ObfsMinPacketSize: 512,
		ObfsMaxPacketSize: 1200,

		TuicCongestionControl: "bbr",
		TuicAuthTimeout:       3,
		TuicHeartbeat:         10,
		TuicUDPRelayMode:      "native",

		HopPort:     "10000-11000",
		HopInterval: "5-10",

		AnyTLSIdleSessionCheckInterval: 30,
		AnyTLSIdleSessionTimeout:       30,
		AnyTLSMinIdleSession:           1,

		TLSType:       "none",
		ALPN:          "h2,http/1.1",
		TLSMinVersion: "1.2",
		TLSMaxVersion: "1.3",
		UTLS:          "chrome",

		RealityServerPort: 443,
	}

	switch cfg.Protocol {
	case "hysteria":
		ib.TLSType = "tls"
		ib.ServerName = cfg.SNI
		ib.ALPN = "h3"
		ib.MasqueradeEnabled = true
		ib.MasqueradeType = "proxy"
		ib.MasqueradeURL = masqueradeURL
		ib.RewriteHost = true

		if cfg.ObfsPassword != "" {
			ib.ObfsType = "salamander"
			ib.ObfsPassword = cfg.ObfsPassword
		}

	case "vless", "vmess":
		switch cfg.Network {
		case "tcp":
			ib.Transport = "raw"
		case "ws":
			ib.Transport = "websocket"
			ib.WsPath = cfg.Path
			ib.WsHost = cfg.Host
		default:
			return ib, fmt.Errorf("节点的传输方式 %s 本程序暂不支持", cfg.Network)
		}

		ib.TLSType = cfg.Security
		if cfg.Fingerprint != "" {
			ib.UTLS = cfg.Fingerprint
		}

		switch cfg.Security {
		case "none":
		case "tls":
			ib.ServerName = cfg.SNI
		case "reality":
			shortIDs, _ := json.Marshal([]string{cfg.ShortID})
			ib.RealityServerName = cfg.SNI
			ib.RealityServer = cfg.SNI
			ib.RealityShortIDs = string(shortIDs)
			ib.RealityPrivateKey = cfg.PrivateKey
			ib.RealityPublicKey = cfg.PublicKey
			if cfg.Protocol == "vless" {
				ib.Flow = "xtls-rprx-vision"
			}
		default:
			return ib, fmt.Errorf("节点的安全方式 %s 本程序暂不支持", cfg.Security)
		}

	default:
		return ib, fmt.Errorf("节点的协议 %s 本程序暂不支持", cfg.Protocol)
	}

	return ib, nil
}

func findCert(domain string) (uint, error) {
	var cert database.Cert
	if database.DB.Where(database.Cert{Domain: domain}).First(&cert).Error != nil {
		return 0, fmt.Errorf("证书库里没有域名 %s 的证书，请先申请", domain)
	}
	if cert.CertPath == "" || cert.KeyPath == "" {
		return 0, fmt.Errorf("域名 %s 的证书还没有签发", domain)
	}
	return cert.ID, nil
}
