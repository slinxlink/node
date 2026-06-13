package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/config"
	"github.com/slinxlink/node/internal/update"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": config.Version})
}

func CheckUpdate(c *gin.Context) {
	result, err := update.Check(config.Version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func Update(c *gin.Context) {
	if err := update.Update(config.Version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
