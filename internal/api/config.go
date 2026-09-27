package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/service"
	"github.com/slinxlink/node/internal/sync"
	"github.com/slinxlink/node/internal/util"
)

func GetConfig(c *gin.Context) {
	var config database.Config
	database.DB.First(&config)
	c.JSON(http.StatusOK, config)
}

func UpdateConfig(c *gin.Context) {
	var req database.Config
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var prev database.Config
	database.DB.First(&prev)

	usedPorts := database.UsedPorts()

	if req.Port != prev.Port {
		if msg := util.ValidatePort(req.Port, usedPorts); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	if !strings.HasPrefix(req.Path, "/") || strings.Count(req.Path, "/") != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径必须以 '/' 开头且只能有一个 '/'"})
		return
	}

	if req.SubPort != prev.SubPort {
		if msg := util.ValidatePort(req.SubPort, usedPorts); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	if !strings.HasPrefix(req.SubPath, "/") || strings.Count(req.SubPath, "/") != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订阅路径必须以 '/' 开头且只能有一个 '/'"})
		return
	}

	if req.LogPath != "" && !strings.HasSuffix(req.LogPath, ".log") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日志路径必须以 .log 结尾"})
		return
	}

	updates := map[string]any{
		"domain":              req.Domain,
		"port":                req.Port,
		"path":                req.Path,
		"sub_enable":          req.SubEnable,
		"sub_port":            req.SubPort,
		"sub_path":            req.SubPath,
		"ruleset_auto_update": req.RulesetAutoUpdate,
		"log_enable":          req.LogEnable,
		"log_path":            req.LogPath,
		"log_level":           req.LogLevel,
		"bbr":                 req.BBR,
		"board_enable":        req.BoardEnable,
	}
	database.DB.Model(&prev).Updates(updates)

	util.InitLog(req.LogPath, req.LogLevel, req.LogEnable)

	if prev.BoardEnable && !req.BoardEnable {
		sync.Stop()
		go core.Default.Apply()
	} else if !prev.BoardEnable && req.BoardEnable {
		sync.Start()
		go core.Default.Apply()
	}

	if req.BBR != prev.BBR {
		service.BBRApply(req.BBR)
	}

	util.Info("[config] 面板配置已更新")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ResetConfig(c *gin.Context) {
	var prev database.Config
	database.DB.First(&prev)

	ipv4, ipv6 := util.GetPublicIPs()
	config := database.Config{
		SecretKey:         util.GenerateString(32),
		Username:          "admin",
		Password:          util.GenerateString(12),
		Port:              util.GeneratePort(),
		Path:              "/" + util.GenerateString(8),
		IPv4:              ipv4,
		IPv6:              ipv6,
		SubEnable:         true,
		SubPath:           "/link",
		SubPort:           2096,
		RulesetAutoUpdate: false,
		LogEnable:         true,
		LogLevel:          "info",
		LogPath:           "data/slinx.log",
		BBR:               true,
		BoardEnable:       false,
		Repo:              "https://github.com/slinxlink/node",
	}
	config.ID = prev.ID

	database.DB.Save(&config)
	util.Info("[config] 面板配置已重置")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
