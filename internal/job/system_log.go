package job

import (
	"time"

	"github.com/seekky/slinx-node/internal/database"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var lastJobNetIO []net.IOCountersStat
var lastJobNetTime time.Time

func SystemLog() {
	go func() {
		// 等到下一个15分钟整点
		now := time.Now()
		next := now.Truncate(15 * time.Minute).Add(15 * time.Minute)
		time.Sleep(next.Sub(now))

		// 先执行一次
		writeSystemLog()

		// 之后每15分钟整点执行
		ticker := time.NewTicker(15 * time.Minute)
		for range ticker.C {
			writeSystemLog()
		}
	}()
}
func writeSystemLog() {
	cpuPercent, _ := cpu.Percent(0, false)
	vmStat, _ := mem.VirtualMemory()
	loadStat, _ := load.Avg()
	currentNetIO, _ := net.IOCounters(false)
	now := time.Now()

	upload := int64(0)
	download := int64(0)
	if lastJobNetIO != nil {
		upload = int64(currentNetIO[0].BytesSent - lastJobNetIO[0].BytesSent)
		download = int64(currentNetIO[0].BytesRecv - lastJobNetIO[0].BytesRecv)
	}
	lastJobNetIO = currentNetIO

	log := &database.SystemLog{
		CPU:      cpuPercent[0],
		RAM:      vmStat.UsedPercent,
		Load:     loadStat.Load15,
		Upload:   upload,
		Download: download,
	}
	database.DB.Create(log)
	database.DB.Where("created_at < ?", now.Add(-24*time.Hour)).Delete(&database.SystemLog{})
}
