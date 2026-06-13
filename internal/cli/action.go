package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const releasesURL = "https://slinx.ink/sh"

func runCmd(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func openDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("data/slinx.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func panelStatus() string {
	out := runCmd("systemctl", "is-active", "slinx")
	if out == "active" {
		return lipgloss.NewStyle().Foreground(successColor).Render("● 运行中")
	}
	return lipgloss.NewStyle().Foreground(errorColor).Render("● 已停止")
}

func coreStatus() string {
	out := runCmd("pgrep", "-x", "sing-box")
	if out != "" {
		return lipgloss.NewStyle().Foreground(successColor).Render("● 运行中")
	}
	return lipgloss.NewStyle().Foreground(errorColor).Render("● 已停止")
}

func autoRestartStatus() string {
	out := runCmd("systemctl", "is-enabled", "slinx")
	if out == "enabled" {
		return lipgloss.NewStyle().Foreground(successColor).Render("● 已开启")
	}
	return lipgloss.NewStyle().Foreground(errorColor).Render("● 已关闭")
}

func startPanel() string {
	out := runCmd("systemctl", "start", "slinx")
	if out == "" {
		return renderStatus("面板", "启动成功", true)
	}
	return renderStatus("面板", out, false)
}

func stopPanel() string {
	out := runCmd("systemctl", "stop", "slinx")
	if out == "" {
		return renderStatus("面板", "已停止", true)
	}
	return renderStatus("面板", out, false)
}

func restartPanel() string {
	out := runCmd("systemctl", "restart", "slinx")
	if out == "" {
		return renderStatus("面板", "重启成功", true)
	}
	return renderStatus("面板", out, false)
}

func toggleAutoRestart() string {
	enabled := runCmd("systemctl", "is-enabled", "slinx")
	if enabled == "enabled" {
		runCmd("systemctl", "disable", "slinx")
		return renderStatus("自动重启", "已关闭", false)
	}
	runCmd("systemctl", "enable", "slinx")
	return renderStatus("自动重启", "已开启", true)
}

func resetUrl() string {
	db, err := openDB()
	if err != nil {
		return "数据库读取失败"
	}

	path := "/" + util.GenerateString(8)
	port := util.GeneratePort()
	ipv4, ipv6 := util.GetPublicIPs()
	if err := db.Model(&database.Config{}).Where("id = 1").Updates(map[string]interface{}{
		"path": path,
		"port": port,
		"ipv4": ipv4,
		"ipv6": ipv6,
	}).Error; err != nil {
		return "保存失败: " + err.Error()
	}

	url := fmt.Sprintf("http://%s:%d%s", ipv4, port, path)

	runCmd("systemctl", "restart", "slinx")

	return renderStatus("面板地址", "重置成功", true) + "\n" + renderInfo("登录信息",
		[]string{"新地址", url},
	)
}

func resetCredentials() string {
	db, err := openDB()
	if err != nil {
		return "数据库读取失败"
	}

	username := "admin"
	password := util.GenerateString(12)
	if err := db.Model(&database.Config{}).Where("id = 1").Updates(map[string]interface{}{
		"username": username,
		"password": password,
	}).Error; err != nil {
		return "保存失败: " + err.Error()
	}

	runCmd("systemctl", "restart", "slinx")

	return renderStatus("用户名&密码", "重置成功", true) + "\n" + renderInfo("登录信息",
		[]string{"用户名", username},
		[]string{"密码", password},
	)
}

func showLoginInfo() string {
	db, err := openDB()
	if err != nil {
		return "数据库读取失败"
	}

	var cfg database.Config
	db.First(&cfg)

	var addr string
	if cfg.Domain != "" {
		addr = fmt.Sprintf("https://%s:%d%s", cfg.Domain, cfg.Port, cfg.Path)
	} else if cfg.IPv4 != "" {
		addr = fmt.Sprintf("http://%s:%d%s", cfg.IPv4, cfg.Port, cfg.Path)
	} else {
		addr = fmt.Sprintf("http://localhost:%d%s", cfg.Port, cfg.Path)
	}

	return renderInfo("登录信息",
		[]string{"访问地址", addr},
		[]string{"用户名", cfg.Username},
		[]string{"密码", cfg.Password},
	)
}

func update() string {
	arch := runCmd("uname", "-m")
	if arch == "x86_64" {
		arch = "amd64"
	} else {
		arch = "arm64"
	}

	url := fmt.Sprintf("%s/slinx_linux_%s", releasesURL, arch)
	out := runCmd("wget", "-O", "/etc/slinx/slinx", url)
	if out != "" {
		return renderStatus("更新", "下载失败", false)
	}
	runCmd("chmod", "+x", "/etc/slinx/slinx")
	runCmd("systemctl", "restart", "slinx")
	return renderStatus("更新", "更新成功", true)
}

func uninstall() string {
	runCmd("systemctl", "stop", "slinx")
	runCmd("systemctl", "disable", "slinx")
	runCmd("rm", "-f", "/etc/systemd/system/slinx.service")
	runCmd("rm", "-rf", "/etc/slinx")
	runCmd("rm", "-f", "/usr/local/bin/slinx")
	runCmd("systemctl", "daemon-reload")
	time.Sleep(5 * time.Second)
	return renderStatus("卸载", "卸载成功", true) + "\n" + renderInfo("提示",
		[]string{"", "5秒后自动退出脚本"},
	)
}

func firstRun() string {
	db, err := openDB()
	if err != nil {
		return ""
	}

	var cfg database.Config
	db.First(&cfg)

	if !cfg.StartedAt.IsZero() {
		return ""
	}

	return showLoginInfo()
}
