package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/sub"
)

func GetSubscriptionPage(c *gin.Context, webFS embed.FS) {
	token := c.Param("token")
	data, err := webFS.ReadFile("web/dist/sub.html")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	info := sub.Info(token)
	pageJSON, _ := json.Marshal(info)
	html := strings.Replace(
		string(data),
		"</head>",
		fmt.Sprintf(`<script>window.__SUB_DATA__=%s</script></head>`, pageJSON),
		1,
	)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func GetSubscription(c *gin.Context) {
	token := c.Param("token")
	result := sub.Sub(token)
	if result == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "无效的订阅链接"})
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(result))
}

func GetClashSubscription(c *gin.Context) {
	token := c.Param("token")
	result := sub.Clash(token)
	if result == "" {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=clash.yaml")
	c.String(http.StatusOK, result)
}

func GetSubscriptionUri(c *gin.Context) {
	var req struct {
		User    database.User    `json:"user"`
		Inbound database.Inbound `json:"inbound"`
	}
	c.ShouldBindJSON(&req)
	result := sub.Uri(req.User, req.Inbound)
	c.JSON(http.StatusOK, gin.H{"uri": result})
}

func GetSubscriptionUrl(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	c.ShouldBindJSON(&req)
	urls := sub.Url(req.Token)
	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
