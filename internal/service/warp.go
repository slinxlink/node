package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
	"gorm.io/gorm"
)

const (
	warpAPIBase   = "https://api.cloudflareclient.com/v0a4005"
	warpUserAgent = "okhttp/3.12.1"
)

// ── 数据结构 ──────────────────────────────────────────────────────────────────

// WarpData 用于生成端点 + 展示账户状态，注册/刷新都返回这个，不落库
type WarpData struct {
	Address       string
	Reserved      string
	PublicKey     string
	PrivateKey    string
	PeerAddress   string
	PeerPort      int
	PeerPublicKey string

	DeviceName    string
	DeviceModel   string
	DeviceEnabled bool
	AccountType   string
	Role          string
	WarpPlusData  int64
	Quota         int64
	Usage         int64
}

type cfRegisterResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Model   string `json:"model"`
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	Account struct {
		License     string `json:"license"`
		AccountType string `json:"account_type"`
		Role        string `json:"role"`
		WarpPlus    bool   `json:"warp_plus"`
		PremiumData int64  `json:"premium_data"`
		Quota       int64  `json:"quota"`
		Usage       int64  `json:"usage"`
	} `json:"account"`
	Config struct {
		ClientID string `json:"client_id"`
		Peers    []struct {
			PublicKey string `json:"public_key"`
			Endpoint  struct {
				Host string `json:"host"`
			} `json:"endpoint"`
		} `json:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `json:"v4"`
				V6 string `json:"v6"`
			} `json:"addresses"`
		} `json:"interface"`
	} `json:"config"`
}

// ── 内部工具函数 ──────────────────────────────────────────────────────────────

// parseReserved base64 -> 3字节 -> int数组JSON
func parseReserved(clientID string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return "", err
	}
	if len(raw) < 3 {
		return "", fmt.Errorf("invalid client_id length: %d", len(raw))
	}
	data, err := json.Marshal([]int{int(raw[0]), int(raw[1]), int(raw[2])})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func warpDoRequest(method, url, token string, body any) (*cfRegisterResponse, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", warpUserAgent)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("cloudflare api error: status=%d body=%s", resp.StatusCode, string(data))
	}

	var result cfRegisterResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func buildWarpData(privateKey, publicKey string, resp *cfRegisterResponse) (*WarpData, error) {
	addressJSON, err := json.Marshal([]string{
		resp.Config.Interface.Addresses.V4 + "/32",
		resp.Config.Interface.Addresses.V6 + "/128",
	})
	if err != nil {
		return nil, err
	}

	reserved, err := parseReserved(resp.Config.ClientID)
	if err != nil {
		return nil, err
	}

	if len(resp.Config.Peers) == 0 {
		return nil, fmt.Errorf("cloudflare response missing peers")
	}
	peer := resp.Config.Peers[0]

	host, portStr, err := net.SplitHostPort(peer.Endpoint.Host)
	if err != nil {
		return nil, err
	}
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	return &WarpData{
		Address:       string(addressJSON),
		Reserved:      reserved,
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		PeerAddress:   host,
		PeerPort:      port,
		PeerPublicKey: peer.PublicKey,

		DeviceName:    resp.Name,
		DeviceModel:   resp.Model,
		DeviceEnabled: resp.Enabled,
		AccountType:   resp.Account.AccountType,
		Role:          resp.Account.Role,
		WarpPlusData:  resp.Account.PremiumData,
		Quota:         resp.Account.Quota,
		Usage:         resp.Account.Usage,
	}, nil
}

// ── 对外方法 ─────────────────────────────────────────────────────────────────

// WarpGet 读取当前 Warp 账号记录（没有则返回 nil）
func WarpGet() (*database.Warp, error) {
	var warp database.Warp
	err := database.DB.First(&warp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &warp, nil
}

// WarpDelete 删除当前 WARP 账号记录
func WarpDelete() error {
	return database.DB.Where("1 = 1").Delete(&database.Warp{}).Error
}

// WarpRegister 调用 Cloudflare 注册新账号 / 换IP，写入 Warp 表
func WarpRegister() (*WarpData, error) {
	privateKey, publicKey, err := util.GenerateWireguardKeyPair()
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	body := map[string]any{
		"key":   publicKey,
		"tos":   time.Now().UTC().Format(time.RFC3339),
		"type":  "PC",
		"name":  hostname,
		"model": "slinx",
	}

	resp, err := warpDoRequest(http.MethodPost, warpAPIBase+"/reg", "", body)
	if err != nil {
		return nil, err
	}

	database.DB.Where("1 = 1").Delete(&database.Warp{})
	warp := database.Warp{
		AccessToken: resp.Token,
		DeviceID:    resp.ID,
		LicenseKey:  resp.Account.License,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		AutoUpdate:  0,
	}
	if err := database.DB.Create(&warp).Error; err != nil {
		return nil, err
	}

	return buildWarpData(privateKey, publicKey, resp)
}

// WarpRefresh 用已有账号信息重新拉取当前状态（密钥不变）
func WarpRefresh() (*WarpData, error) {
	warp, err := WarpGet()
	if err != nil {
		return nil, err
	}
	if warp == nil {
		return nil, fmt.Errorf("没有可用的 WARP 账号")
	}

	resp, err := warpDoRequest(http.MethodGet, fmt.Sprintf("%s/reg/%s", warpAPIBase, warp.DeviceID), warp.AccessToken, nil)
	if err != nil {
		return nil, err
	}

	if resp.Account.License != "" && resp.Account.License != warp.LicenseKey {
		warp.LicenseKey = resp.Account.License
		database.DB.Save(warp)
	}

	return buildWarpData(warp.PrivateKey, warp.PublicKey, resp)
}

func WarpSetAutoUpdate(day int) (*database.Warp, error) {
	warp, err := WarpGet()
	if err != nil {
		return nil, err
	}
	if warp == nil {
		return nil, fmt.Errorf("没有可用的 WARP 账号")
	}
	warp.AutoUpdate = day
	if err := database.DB.Save(warp).Error; err != nil {
		return nil, err
	}
	return warp, nil
}

// WarpSetLicense 绑定/升级 WARP+ 许可证密钥
func WarpSetLicense(license string) (*WarpData, error) {
	warp, err := WarpGet()
	if err != nil {
		return nil, err
	}
	if warp == nil {
		return nil, fmt.Errorf("没有可用的 WARP 账号")
	}

	resp, err := warpDoRequest(
		http.MethodPut,
		fmt.Sprintf("%s/reg/%s/account", warpAPIBase, warp.DeviceID),
		warp.AccessToken,
		map[string]string{"license": license},
	)
	if err != nil {
		return nil, err
	}

	warp.LicenseKey = resp.Account.License
	if err := database.DB.Save(warp).Error; err != nil {
		return nil, err
	}

	return buildWarpData(warp.PrivateKey, "", resp)
}
