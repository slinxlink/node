package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/config"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/service"
	"github.com/slinxlink/node/internal/sync"
	"github.com/slinxlink/node/internal/util"
)

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Config)
}

func UpdateConfig(c *gin.Context) {
	var cfg database.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	usedPorts := database.UsedPorts()

	if cfg.Port != config.Config.Port {
		if msg := util.ValidatePort(cfg.Port, usedPorts); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	if !strings.HasPrefix(cfg.Path, "/") || strings.Count(cfg.Path, "/") != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径必须以 '/' 开头且只能有一个 '/'"})
		return
	}

	if cfg.SubPort != config.Config.SubPort {
		if msg := util.ValidatePort(cfg.SubPort, usedPorts); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	if !strings.HasPrefix(cfg.SubPath, "/") || strings.Count(cfg.SubPath, "/") != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订阅路径必须以 '/' 开头且只能有一个 '/'"})
		return
	}

	if cfg.LogPath != "" && !strings.HasSuffix(cfg.LogPath, ".log") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日志路径必须以 .log 结尾"})
		return
	}

	prev := config.Config

	cfg.ID = config.Config.ID
	database.DB.Save(&cfg)

	config.Config = cfg
	util.InitLog(cfg.LogPath, cfg.LogLevel, cfg.LogEnable)

	if prev.BoardEnable && !cfg.BoardEnable {
		sync.Stop()
		go core.Default.Apply()
	} else if !prev.BoardEnable && cfg.BoardEnable {
		sync.Start()
		go core.Default.Apply()
	}

	if cfg.BBR != prev.BBR {
		service.BBRApply(cfg.BBR)
	}

	util.Info("[config] 面板配置已更新")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ResetConfig(c *gin.Context) {
	ipv4, ipv6 := util.GetPublicIPs()
	cfg := database.Config{
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
	cfg.ID = config.Config.ID

	database.DB.Save(&cfg)
	config.Config = cfg

	util.Info("[config] 面板配置已重置")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
