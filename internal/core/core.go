package core

import (
	"strings"

	"github.com/slinxlink/node/internal/database"
)

var validLogLevels = map[string]bool{
	"trace": true, "debug": true, "info": true,
	"warn": true, "error": true, "fatal": true, "panic": true,
}

func Update(cfg database.Core) string {
	if cfg.LogLevel != "" && !validLogLevels[cfg.LogLevel] {
		return "日志等级不合法"
	}

	if cfg.LogPath != "" && !strings.HasSuffix(cfg.LogPath, ".log") {
		return "日志路径必须以 .log 结尾"
	}

	database.DB.Save(&cfg)

	if err := Default.Apply(); err != nil {
		return err.Error()
	}

	return ""
}

func Reset() {
	database.DB.Model(&database.Core{}).Where("id = 1").Updates(map[string]interface{}{
		"log_enable": true,
		"log_level":  "info",
		"log_path":   "data/sing-box.log",
	})
	go Default.Apply()
}
