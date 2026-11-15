package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/slinxlink/node/internal/database"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchIPInfo(source string) (*database.IP, error) {
	switch source {
	case "ip.net.coffee":
		return fetchFromNetCoffee()
	case "ip-api.com":
		return fetchFromIpApiCom()
	case "ippure.com":
		return fetchFromIppure()
	case "ip2location.io":
		return fetchFromIp2Location()
	default:
		return nil, fmt.Errorf("不支持的数据源: %s", source)
	}
}

// ────────────────────────────────────────────────
// ip.net.coffee
// ────────────────────────────────────────────────

// 走 IPv4 的客户端：这里只记 IPv4 的结果，机器有 IPv6 时不能让它自己挑出口
var ipv4Client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp4", addr)
		},
	},
}

type netCoffeeResponse struct {
	IP                    string `json:"ip"`
	ASN                   int    `json:"asn"`
	ASOrganization        string `json:"asOrganization"`
	CountryCode           string `json:"countryCode"`
	City                  string `json:"city"`
	RegisteredCountryCode string `json:"registered_country_code"`
	GeoSources            []struct {
		Src string   `json:"src"`
		Lat *float64 `json:"lat"`
		Lon *float64 `json:"lon"`
	} `json:"geo_sources"`
	IsPublicService bool `json:"is_public_service"`
	IsDatacenter    bool `json:"is_datacenter"`
	IsResidential   bool `json:"isResidential"`
	IsMobile        bool `json:"is_mobile"`
	IsVPN           bool `json:"is_vpn"`
	IsProxy         bool `json:"is_proxy"`
	IsTor           bool `json:"is_tor"`
}

// 接口要带上 IP，先问它自己的 trace 拿到本机出口，再查这个 IP
func fetchFromNetCoffee() (*database.IP, error) {
	trace, err := ipv4Client.Get("https://ip.net.coffee/cdn-cgi/trace")
	if err != nil {
		return nil, err
	}
	defer trace.Body.Close()

	ip := ""
	scanner := bufio.NewScanner(trace.Body)
	for scanner.Scan() {
		if value, ok := strings.CutPrefix(scanner.Text(), "ip="); ok {
			ip = strings.TrimSpace(value)
		}
	}
	if ip == "" {
		return nil, fmt.Errorf("ip.net.coffee 没有返回出口 IP")
	}

	resp, err := httpClient.Get("https://ip.net.coffee/api/ip/lookup/" + ip)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ip.net.coffee 查询失败: %d", resp.StatusCode)
	}

	var result netCoffeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.IP == "" {
		return nil, fmt.Errorf("ip.net.coffee 查询失败")
	}

	// IP 类型：代理类(Tor、VPN、代理)优先，其次机房、移动、家宽，都不是就判断不了
	ipType := ""
	switch {
	case result.IsTor || result.IsVPN || result.IsProxy:
		ipType = "proxy"
	case result.IsDatacenter:
		ipType = "hosting"
	case result.IsMobile:
		ipType = "mobile"
	case result.IsResidential:
		ipType = "isp"
	}

	// 原生性：使用地和注册地一致是原生，不一致是广播；公共服务的 IP 或缺一个国家码就不判断
	native := ""
	if !result.IsPublicService && result.CountryCode != "" && result.RegisteredCountryCode != "" {
		native = "broadcast"
		if strings.EqualFold(result.CountryCode, result.RegisteredCountryCode) {
			native = "native"
		}
	}

	// 经纬度：多个来源各给一组，g3 最细(精确到区)，没有就用 g2，再没有用 g7
	lat, lon := pickGeo(result.GeoSources)

	record := &database.IP{
		Source:    "ip.net.coffee",
		IPVersion: "v4",
		IP:        result.IP,
		ASN:       fmt.Sprintf("AS%d", result.ASN),
		ASNOrg:    result.ASOrganization,
		Country:   strings.ToUpper(result.CountryCode),
		City:      result.City,
		IPType:    ipType,
		Native:    native,
		Latitude:  lat,
		Longitude: lon,
		UpdatedAt: time.Now(),
	}

	return upsertIP(record)
}

// 按 g3、g2、g7 的顺序取第一个经纬度都有值的来源，都没有返回 0, 0
func pickGeo(sources []struct {
	Src string   `json:"src"`
	Lat *float64 `json:"lat"`
	Lon *float64 `json:"lon"`
}) (float64, float64) {
	for _, name := range []string{"g3", "g2", "g7"} {
		for _, item := range sources {
			if item.Src == name && item.Lat != nil && item.Lon != nil {
				return *item.Lat, *item.Lon
			}
		}
	}

	return 0, 0
}

// ────────────────────────────────────────────────
// ip-api.com
// ────────────────────────────────────────────────

type ipApiComResponse struct {
	Status      string  `json:"status"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	ISP         string  `json:"isp"`
	AS          string  `json:"as"`
	Mobile      bool    `json:"mobile"`
	Proxy       bool    `json:"proxy"`
	Hosting     bool    `json:"hosting"`
	Query       string  `json:"query"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func fetchFromIpApiCom() (*database.IP, error) {
	resp, err := httpClient.Get("http://ip-api.com/json/?fields=status,message,countryCode,city,lat,lon,isp,as,mobile,proxy,hosting,query")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ipApiComResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("ip-api.com 查询失败")
	}

	// 机房、移动、代理都不是就算家宽
	ipType := "isp"
	if result.Hosting {
		ipType = "hosting"
	} else if result.Mobile {
		ipType = "mobile"
	} else if result.Proxy {
		ipType = "proxy"
	}

	record := &database.IP{
		Source:    "ip-api.com",
		IPVersion: "v4",
		IP:        result.Query,
		ASN:       result.AS,
		ASNOrg:    result.ISP,
		Country:   strings.ToUpper(result.CountryCode),
		City:      result.City,
		IPType:    ipType,
		Latitude:  result.Lat,
		Longitude: result.Lon,
		UpdatedAt: time.Now(),
	}

	return upsertIP(record)
}

// ────────────────────────────────────────────────
// ippure.com
// ────────────────────────────────────────────────

type ippureResponse struct {
	IP             string `json:"ip"`
	ASN            int    `json:"asn"`
	AsOrganization string `json:"asOrganization"`
	Country        string `json:"country"`
	CountryCode    string `json:"countryCode"`
	City           string `json:"city"`
	FraudScore     int    `json:"fraudScore"`
	IsResidential  bool   `json:"isResidential"`
	IsBroadcast    *bool  `json:"isBroadcast"` // 没返回就是 nil，不能当成原生
	Latitude       string `json:"latitude"`    // 接口给的是字符串
	Longitude      string `json:"longitude"`
}

func fetchFromIppure() (*database.IP, error) {
	resp, err := httpClient.Get("https://my.ippure.com/v1/info")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ippureResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 不是家宽就算机房
	ipType := "hosting"
	if result.IsResidential {
		ipType = "isp"
	}

	native := ""
	if result.IsBroadcast != nil {
		native = "native"
		if *result.IsBroadcast {
			native = "broadcast"
		}
	}

	// 经纬度是字符串，转不了就当没有
	lat, _ := strconv.ParseFloat(result.Latitude, 64)
	lon, _ := strconv.ParseFloat(result.Longitude, 64)

	record := &database.IP{
		Source:    "ippure.com",
		IPVersion: "v4",
		IP:        result.IP,
		ASN:       fmt.Sprintf("AS%d", result.ASN),
		ASNOrg:    result.AsOrganization,
		Country:   strings.ToUpper(result.CountryCode),
		City:      result.City,
		IPType:    ipType,
		Native:    native,
		Latitude:  lat,
		Longitude: lon,
		UpdatedAt: time.Now(),
	}

	return upsertIP(record)
}

// ────────────────────────────────────────────────
// ip2location.io
// ────────────────────────────────────────────────

type ip2LocationResponse struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CityName    string  `json:"city_name"`
	ASN         string  `json:"asn"`
	As          string  `json:"as"`
	IsProxy     bool    `json:"is_proxy"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func fetchFromIp2Location() (*database.IP, error) {
	resp, err := httpClient.Get("https://api.ip2location.io/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ip2LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ipType := ""
	if result.IsProxy {
		ipType = "proxy"
	}

	record := &database.IP{
		Source:    "ip2location.io",
		IPVersion: "v4",
		IP:        result.IP,
		ASN:       fmt.Sprintf("AS%s", result.ASN),
		ASNOrg:    result.As,
		Country:   strings.ToUpper(result.CountryCode),
		City:      result.CityName,
		IPType:    ipType,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		UpdatedAt: time.Now(),
	}

	return upsertIP(record)
}

// ────────────────────────────────────────────────
// 通用工具
// ────────────────────────────────────────────────

func upsertIP(record *database.IP) (*database.IP, error) {
	var existing database.IP
	database.DB.Where(database.IP{
		Source:    record.Source,
		IPVersion: record.IPVersion,
	}).First(&existing)

	if existing.ID == 0 {
		database.DB.Create(record)
	} else {
		record.ID = existing.ID
		database.DB.Save(record)
	}

	return record, nil
}
