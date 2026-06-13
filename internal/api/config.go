package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seekky/slinx-node/internal/config"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/util"
)

func GetConfig(c *gin.Context) {
	var config database.Config
	database.DB.First(&config)
	c.JSON(http.StatusOK, config)
}

func UpdateConfig(c *gin.Context) {
	var cfg database.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if msg := config.Update(cfg); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	util.Info("[config] 面板配置已更新")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ApplyCertCF(c *gin.Context)   { c.JSON(200, gin.H{"message": "ok"}) }
func ApplyCertHTTP(c *gin.Context) { c.JSON(200, gin.H{"message": "ok"}) }
