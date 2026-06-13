package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/slinxlink/node/internal/database"
)

var lastNetIO []net.IOCountersStat
var lastNetTime time.Time

func GetSystemStatus(c *gin.Context) {
	cpuPercent, _ := cpu.Percent(0, false)
	cpuInfo, _ := cpu.Info()
	cores, _ := cpu.Counts(true)
	vmStat, _ := mem.VirtualMemory()
	swapStat, _ := mem.SwapMemory()
	diskStat, _ := disk.Usage("/")
	tcpConn, _ := net.Connections("tcp")
	udpConn, _ := net.Connections("udp")
	currentNetIO, _ := net.IOCounters(false)
	now := time.Now()

	mhz := 0.0
	if len(cpuInfo) > 0 {
		mhz = cpuInfo[0].Mhz
	}

	percent := 0.0
	if len(cpuPercent) > 0 {
		percent = cpuPercent[0]
	}

	uploadSpeed := 0.0
	downloadSpeed := 0.0
	if lastNetIO != nil {
		elapsed := now.Sub(lastNetTime).Seconds()
		uploadSpeed = float64(currentNetIO[0].BytesSent-lastNetIO[0].BytesSent) / elapsed
		downloadSpeed = float64(currentNetIO[0].BytesRecv-lastNetIO[0].BytesRecv) / elapsed
	}
	lastNetIO = currentNetIO
	lastNetTime = now

	c.JSON(200, gin.H{
		"cpu": gin.H{
			"percent": percent,
			"cores":   cores,
			"mhz":     mhz,
		},
		"ram": gin.H{
			"percent": vmStat.UsedPercent,
			"used":    vmStat.Used,
			"total":   vmStat.Total,
		},
		"swap": gin.H{
			"percent": swapStat.UsedPercent,
			"used":    swapStat.Used,
			"total":   swapStat.Total,
		},
		"disk": gin.H{
			"percent": diskStat.UsedPercent,
			"used":    diskStat.Used,
			"total":   diskStat.Total,
		},
		"connections": gin.H{
			"tcp": len(tcpConn),
			"udp": len(udpConn),
		},
		"network": gin.H{
			"upload":   uploadSpeed,
			"download": downloadSpeed,
		},
	})
}

func GetStats(c *gin.Context) {
	var stats database.Stats
	database.DB.Where("inbound_id = 0").First(&stats)
	c.JSON(200, stats)
}

func GetSystemLog(c *gin.Context) {
	var logs []database.SystemLog
	database.DB.Order("created_at asc").Find(&logs)
	c.JSON(200, logs)
}
