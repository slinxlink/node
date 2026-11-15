package job

import (
	"time"

	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/util"
)

func TrafficReset() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			resetTraffic()
		}
	}()
}

func resetTraffic() {
	defer func() {
		if r := recover(); r != nil {
			util.Error("[stats] 重置流量时出错: %v", r)
		}
	}()

	core.ResetTraffic()
	core.CheckLimits()
}
