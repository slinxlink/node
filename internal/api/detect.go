package api

import (
	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/service"
	"github.com/slinxlink/node/internal/util"
)

// GET /detect/ip 读库返回
func DetectIP(c *gin.Context) {
	var records []database.IP
	database.DB.Find(&records)
	c.JSON(200, gin.H{"data": records})
}

// POST /detect/ip/fetch 触发查询
func FetchIP(c *gin.Context) {
	var req struct {
		Source string `json:"source"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Source == "" {
		req.Source = "ipapi.is"
	}

	util.Info("[detect] 查询 IP 信息: %s", req.Source)
	record, err := service.FetchIPInfo(req.Source)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": record})
}

// GET /detect/unlock 读库返回
func DetectUnlock(c *gin.Context) {
	var records []database.Unlock
	database.DB.Find(&records)
	c.JSON(200, gin.H{"data": records})
}

// POST /detect/unlock/fetch 触发查询
func FetchUnlock(c *gin.Context) {
	util.Info("[detect] 查询流媒体解锁")
	records, err := service.FetchUnlockInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": records})
}

// GET /detect/back-route 读库返回
func DetectBackRoute(c *gin.Context) {
	records, err := service.GetBackRouteInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": records})
}

// POST /detect/back-route/fetch 触发查询
func FetchBackRoute(c *gin.Context) {
	util.Info("[detect] 查询回程路由")
	records, err := service.FetchBackRouteInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": records})
}
