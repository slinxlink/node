package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/util"
)

func GeneratePort(c *gin.Context) {
	var cfg database.Config
	database.DB.First(&cfg)

	for {
		port := util.GeneratePort()

		// 检查保留端口和面板端口
		if msg := util.ValidatePort(port, cfg.Port, []int{}); msg != "" {
			continue
		}

		// 检查是否被入站占用
		var count int64
		database.DB.Model(&database.Inbound{}).Where("port = ?", port).Count(&count)
		if count == 0 {
			c.JSON(http.StatusOK, gin.H{"port": port})
			return
		}
	}
}

func GenerateRealityTarget(c *gin.Context) {
	domain := util.GenerateRealityTarget()
	c.JSON(http.StatusOK, gin.H{"domain": domain})
}

func GenerateRealityKeyPair(c *gin.Context) {
	priv, pub, err := util.GenerateRealityKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"private_key": priv, "public_key": pub})
}

func GenerateShortIDs(c *gin.Context) {
	ids, err := util.GenerateShortIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"short_ids": ids})
}

func GenerateToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": util.GenerateString(16)})
}

func GenerateUUID(c *gin.Context) {
	uuid, err := util.GenerateUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uuid": uuid})
}

func GeneratePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"password": util.GenerateString(16)})
}
