package config

import (
	"strings"

	"github.com/seekky/slinx-node/internal/core"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/sync"
	"github.com/seekky/slinx-node/internal/util"
)

var Config database.Config

func Load() error {
	return database.DB.First(&Config).Error
}

func Update(cfg database.Config) string {
	if msg := util.ValidatePort(cfg.Port, 0, nil); msg != "" {
		return msg
	}

	if !strings.HasPrefix(cfg.Path, "/") || strings.Count(cfg.Path, "/") != 1 {
		return "路径必须以 '/' 开头且只能有一个 '/'"
	}

	if cfg.LogPath != "" && !strings.HasSuffix(cfg.LogPath, ".log") {
		return "日志路径必须以 .log 结尾"
	}

	prev := Config.BoardEnable
	database.DB.Save(&cfg)
	Config = cfg
	util.InitLog(cfg.LogPath, cfg.LogLevel, cfg.LogEnable)

	if prev && !cfg.BoardEnable {
		sync.Stop()
		go core.Default.Apply()
	} else if !prev && cfg.BoardEnable {
		sync.Start()
		go core.Default.Apply()
	}

	return ""
}
