package api

import (
	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/service"
	"github.com/slinxlink/node/internal/util"
)

func GetWarp(c *gin.Context) {
	warp, err := service.WarpGet()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": warp})
}

func DeleteWarp(c *gin.Context) {
	util.Info("[warp] 删除账号")
	if err := service.WarpDelete(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}

func RegisterWarp(c *gin.Context) {
	util.Info("[warp] 注册账号")
	data, err := service.WarpRegister()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func RefreshWarp(c *gin.Context) {
	util.Info("[warp] 刷新账号")
	data, err := service.WarpRefresh()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func SetWarpAutoUpdate(c *gin.Context) {
	var day int
	if err := c.ShouldBindJSON(&day); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	warp, err := service.WarpSetAutoUpdate(day)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": warp})
}

func SetWarpLicense(c *gin.Context) {
	var license string
	if err := c.ShouldBindJSON(&license); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	util.Info("[warp] 修改许可证")
	data, err := service.WarpSetLicense(license)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}
