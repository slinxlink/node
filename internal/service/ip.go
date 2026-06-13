package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/slinxlink/node/internal/database"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchIPInfo(source string) (*database.IP, error) {
	switch source {
	case "ipapi.is":
		return fetchFromIpapiIs()
	case "ip-api.com":
		return fetchFromIpApiCom()
	case "ippure.com":
		return fetchFromIppure()
	default:
		return nil, fmt.Errorf("不支持的数据源: %s", source)
	}
}

// ────────────────────────────────────────────────
// ipapi.is
// ────────────────────────────────────────────────

type ipapiIsResponse struct {
	IP  string `json:"ip"`
	ASN struct {
		ASN     int    `json:"asn"`
		Descr   string `json:"descr"`
		Org     string `json:"org"`
		Type    string `json:"type"`
		Country string `json:"country"`
	} `json:"asn"`
	Company struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"company"`
	Location struct {
		CountryCode string `json:"country_code"`
		City        string `json:"city"`
	} `json:"location"`
	IsMobile bool `json:"is_mobile"`
}

func fetchFromIpapiIs() (*database.IP, error) {
	resp, err := httpClient.Get("https://api.ipapi.is")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ipapiIsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ipType := result.Company.Type
	if result.IsMobile {
		ipType = "mobile"
	} else if ipType == "" {
		ipType = result.ASN.Type
	}

	record := &database.IP{
		Source:          "ipapi.is",
		IPVersion:       "v4",
		IP:              result.IP,
		ASN:             result.ASN.Descr,
		ASNOrg:          result.ASN.Org,
		RegisterCountry: strings.ToUpper(result.ASN.Country),
		Country:         strings.ToUpper(result.Location.CountryCode),
		City:            result.Location.City,
		OrgType:         result.Company.Type,
		IPType:          ipType,
		UpdatedAt:       time.Now(),
	}

	return upsertIP(record)
}

// ────────────────────────────────────────────────
// ip-api.com
// ────────────────────────────────────────────────

type ipApiComResponse struct {
	Status      string `json:"status"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	AS          string `json:"as"`
	Mobile      bool   `json:"mobile"`
	Proxy       bool   `json:"proxy"`
	Hosting     bool   `json:"hosting"`
	Query       string `json:"query"`
}

func fetchFromIpApiCom() (*database.IP, error) {
	resp, err := httpClient.Get("http://ip-api.com/json/?fields=status,message,countryCode,city,isp,as,mobile,proxy,hosting,query")
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

	ipType := "business"
	if result.Hosting {
		ipType = "hosting"
	} else if result.Mobile {
		ipType = "mobile"
	} else if result.Proxy {
		ipType = "proxy"
	} else {
		ipType = "isp"
	}

	record := &database.IP{
		Source:          "ip-api.com",
		IPVersion:       "v4",
		IP:              result.Query,
		ASN:             result.AS,
		ASNOrg:          result.ISP,
		RegisterCountry: "",
		Country:         strings.ToUpper(result.CountryCode),
		City:            result.City,
		OrgType:         ipType,
		IPType:          ipType,
		UpdatedAt:       time.Now(),
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

	ipType := "hosting"
	if result.IsResidential {
		ipType = "isp"
	}

	record := &database.IP{
		Source:          "ippure.com",
		IPVersion:       "v4",
		IP:              result.IP,
		ASN:             fmt.Sprintf("AS%d", result.ASN),
		ASNOrg:          result.AsOrganization,
		RegisterCountry: "",
		Country:         strings.ToUpper(result.CountryCode),
		City:            result.City,
		OrgType:         ipType,
		IPType:          ipType,
		UpdatedAt:       time.Now(),
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
