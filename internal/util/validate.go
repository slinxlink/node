package util

import "strings"

var reservedPorts = map[int]bool{
	21:   true, // FTP 文件传输
	22:   true, // SSH 远程登录
	23:   true, // Telnet 远程登录
	25:   true, // SMTP 邮件发送
	53:   true, // DNS 域名解析
	80:   true, // HTTP
	110:  true, // POP3 邮件接收
	143:  true, // IMAP 邮件接收
	465:  true, // SMTPS 加密邮件发送
	587:  true, // SMTP 邮件提交
	993:  true, // IMAPS 加密邮件接收
	995:  true, // POP3S 加密邮件接收
	2048: true, // sing-box 流量统计接口
	3306: true, // MySQL 数据库
	3389: true, // RDP Windows 远程桌面
	5432: true, // PostgreSQL 数据库
	6379: true, // Redis 缓存数据库
	9090: true, // sing-box Clash API
}

func ValidatePort(port int, panelPort int, usedPorts []int) string {
	if reservedPorts[port] {
		return "不能使用系统保留端口"
	}
	if panelPort != 0 && port == panelPort {
		return "不能使用面板端口"
	}
	for _, p := range usedPorts {
		if port == p {
			return "端口已被入站协议占用"
		}
	}
	return ""
}

func ValidateDomain(domain string) bool {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false
	}
	if parts[0] == "" {
		return false
	}
	tld := parts[len(parts)-1]
	if tld == "" || len(tld) < 2 {
		return false
	}
	return true
}
