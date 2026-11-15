package job

import (
	"time"

	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/util"
)

func Traffic() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			collectTraffic()
		}
	}()
}

func collectTraffic() {
	defer func() {
		if r := recover(); r != nil {
			util.Error("[stats] 采集流量时出错: %v", r)
		}
	}()

	if core.Default.Status() == "stopped" {
		return
	}

	if err := core.CollectTraffic(); err != nil {
		util.Error("[stats] 采集流量失败: %v", err)
	}

	if err := core.CollectOnline(); err != nil {
		util.Error("[stats] 获取在线用户失败: %v", err)
	}

	core.CheckLimits()
}
