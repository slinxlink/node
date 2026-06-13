package service

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/seekky/slinx-node/internal/database"
)

// 每个目标：城市、运营商、3个备用IP
type routeTarget struct {
	City string
	ISP  string
	IPs  []string
}

var routeTargets = []routeTarget{
	{City: "shanghai", ISP: "telecom", IPs: []string{"202.96.209.5", "61.129.7.1", "180.169.255.1"}},
	{City: "shanghai", ISP: "unicom", IPs: []string{"210.22.97.1", "58.247.0.1", "116.228.111.1"}},
	{City: "shanghai", ISP: "mobile", IPs: []string{"211.136.112.50", "117.131.0.1", "120.196.0.1"}},
	{City: "beijing", ISP: "telecom", IPs: []string{"219.141.140.10", "202.106.0.20", "61.51.0.1"}},
	{City: "beijing", ISP: "unicom", IPs: []string{"123.123.123.123", "202.106.50.1", "60.247.0.1"}},
	{City: "beijing", ISP: "mobile", IPs: []string{"221.130.33.52", "211.137.130.1", "120.197.0.1"}},
	{City: "guangzhou", ISP: "telecom", IPs: []string{"119.145.0.1", "113.108.0.1", "58.60.0.1"}},
	{City: "guangzhou", ISP: "unicom", IPs: []string{"210.21.4.1", "120.80.0.1", "183.56.0.1"}},
	{City: "guangzhou", ISP: "mobile", IPs: []string{"120.196.165.24", "211.139.0.1", "117.135.0.1"}},
}

// mtr JSON 结构
type mtrReport struct {
	Report struct {
		Hubs []struct {
			Host string  `json:"host"`
			Loss float64 `json:"Loss%"`
		} `json:"hubs"`
	} `json:"report"`
}

// 运行 mtr，失败自动尝试下一个备用 IP
func runMTR(ips []string) ([]string, error) {
	for _, ip := range ips {
		out, err := exec.Command("mtr", "--tcp", "--port", "80",
			"--report", "--report-cycles", "5", "--json", ip).Output()
		if err != nil {
			continue
		}
		var report mtrReport
		if err := json.Unmarshal(out, &report); err != nil {
			continue
		}
		var hosts []string
		for _, hub := range report.Report.Hubs {
			if hub.Host != "???" && hub.Loss < 100 {
				hosts = append(hosts, hub.Host)
			}
		}
		if len(hosts) > 0 {
			return hosts, nil
		}
	}
	return nil, fmt.Errorf("all IPs failed")
}

// IP段匹配
func ipHasPrefix(ip, prefix string) bool {
	parts := strings.Split(prefix, ".")
	ipParts := strings.Split(ip, ".")
	for i, p := range parts {
		if p == "*" {
			continue
		}
		if i >= len(ipParts) || ipParts[i] != p {
			return false
		}
	}
	return true
}

func isPublicIP(host string) bool {
	ip := net.ParseIP(host)
	return ip != nil && !ip.IsPrivate() && !ip.IsLoopback()
}

func detectTelecom(hosts []string) string {
	cn2Count := 0
	for _, h := range hosts {
		if !isPublicIP(h) {
			continue
		}
		if ipHasPrefix(h, "59.43.*.*") {
			cn2Count++
		}
	}
	if cn2Count >= 2 {
		return "CN2 GIA"
	}
	if cn2Count == 1 {
		return "CN2 GT"
	}
	return "163"
}

func detectUnicom(hosts []string) string {
	for _, h := range hosts {
		if !isPublicIP(h) {
			continue
		}
		if ipHasPrefix(h, "210.51.*.*") ||
			ipHasPrefix(h, "218.105.*.*") {
			return "9929"
		}
	}
	for _, h := range hosts {
		if !isPublicIP(h) {
			continue
		}
		if ipHasPrefix(h, "219.158.*.*") {
			return "4837"
		}
	}
	return "169"
}

func detectMobile(hosts []string) string {
	for _, h := range hosts {
		if !isPublicIP(h) {
			continue
		}
		if ipHasPrefix(h, "223.120.*.*") ||
			ipHasPrefix(h, "223.119.*.*") ||
			ipHasPrefix(h, "223.118.*.*") {
			return "CMIN2"
		}
	}
	for _, h := range hosts {
		if !isPublicIP(h) {
			continue
		}
		if ipHasPrefix(h, "221.183.*.*") ||
			ipHasPrefix(h, "211.136.*.*") ||
			ipHasPrefix(h, "211.137.*.*") {
			return "CMI Basic"
		}
	}
	return "CMI"
}

func detectLine(isp string, hosts []string) string {
	switch isp {
	case "telecom":
		return detectTelecom(hosts)
	case "unicom":
		return detectUnicom(hosts)
	case "mobile":
		return detectMobile(hosts)
	}
	return "未知"
}

type routeResult struct {
	City string
	ISP  string
	Line string
}

func FetchRouteInfo() ([]database.Route, error) {
	results := make([]routeResult, len(routeTargets))
	var wg sync.WaitGroup

	for i, target := range routeTargets {
		wg.Add(1)
		go func(idx int, t routeTarget) {
			defer wg.Done()
			hosts, err := runMTR(t.IPs)
			line := "检测失败"
			if err == nil {
				line = detectLine(t.ISP, hosts)
			}
			results[idx] = routeResult{City: t.City, ISP: t.ISP, Line: line}
		}(i, target)
	}

	wg.Wait()

	// 按城市汇总
	cityMap := map[string]*database.Route{}
	for _, r := range results {
		if _, ok := cityMap[r.City]; !ok {
			cityMap[r.City] = &database.Route{City: r.City, UpdatedAt: time.Now()}
		}
		switch r.ISP {
		case "telecom":
			cityMap[r.City].Telecom = r.Line
		case "unicom":
			cityMap[r.City].Unicom = r.Line
		case "mobile":
			cityMap[r.City].Mobile = r.Line
		}
	}

	var routes []database.Route
	for _, city := range []string{"shanghai", "beijing", "guangzhou"} {
		if r, ok := cityMap[city]; ok {
			var existing database.Route
			database.DB.Where(database.Route{City: city}).First(&existing)
			if existing.ID == 0 {
				database.DB.Create(r)
			} else {
				r.ID = existing.ID
				database.DB.Save(r)
			}
			routes = append(routes, *r)
		}
	}

	return routes, nil
}

func GetRouteInfo() ([]database.Route, error) {
	var routes []database.Route
	database.DB.Order("id asc").Find(&routes)
	return routes, nil
}
