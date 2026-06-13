package job

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/slinxlink/node/internal/database"
	"gorm.io/gorm"
)

var lastStatsNetIO []net.IOCountersStat

func Stats() {
	go func() {
		now := time.Now()
		next := now.Truncate(15 * time.Minute).Add(15 * time.Minute)
		time.Sleep(next.Sub(now))

		updateStats()

		ticker := time.NewTicker(15 * time.Minute)
		for range ticker.C {
			updateStats()
		}
	}()
}

func updateStats() {
	currentNetIO, _ := net.IOCounters(false)

	if lastStatsNetIO != nil {
		upload := int64(currentNetIO[0].BytesSent - lastStatsNetIO[0].BytesSent)
		download := int64(currentNetIO[0].BytesRecv - lastStatsNetIO[0].BytesRecv)

		database.DB.Model(&database.Stats{}).Where("inbound_id = 0").Updates(map[string]interface{}{
			"upload":   gorm.Expr("upload + ?", upload),
			"download": gorm.Expr("download + ?", download),
		})
	}
	lastStatsNetIO = currentNetIO
}
