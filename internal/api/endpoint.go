package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/route"
	"github.com/slinxlink/node/internal/util"
)

func GetEndpoints(c *gin.Context) {
	var list []database.Endpoint
	database.DB.Find(&list)
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func SaveEndpoint(c *gin.Context) {
	var endpoint database.Endpoint
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endpoint.Tag = strings.TrimSpace(endpoint.Tag)

	if msg := util.ValidateTag(endpoint.Tag); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var existing database.Endpoint
	if err := database.DB.Where("tag = ? AND id != ?", endpoint.Tag, endpoint.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该标签已存在"})
		return
	}

	if endpoint.Type == "wireguard" {
		var address []string
		if err := json.Unmarshal([]byte(endpoint.Address), &address); err != nil || len(address) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "本地地址格式错误"})
			return
		}
		if !util.ValidateIPv4CIDR(address[0]) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IPv4 地址段格式错误"})
			return
		}
		if !util.ValidateIPv6CIDR(address[1]) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IPv6 地址段格式错误"})
			return
		}
		if endpoint.MTU < 576 || endpoint.MTU > 9000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MTU 必须在 576 - 9000 之间"})
			return
		}
		var reserved []int
		if err := json.Unmarshal([]byte(endpoint.Reserved), &reserved); err != nil || len(reserved) != 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reserved 格式错误，需为3个数字"})
			return
		}
		for _, v := range reserved {
			if v < 0 || v > 255 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Reserved 数值必须在 0-255 之间"})
				return
			}
		}
		if endpoint.PrivateKey == "" || endpoint.PublicKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写密钥对"})
			return
		}
		if endpoint.PeerAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写对端地址"})
			return
		}
		if endpoint.PeerPort < 1 || endpoint.PeerPort > 65535 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "对端端口必须在 1 - 65535 之间"})
			return
		}
		if endpoint.PeerPublicKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写对端公钥"})
			return
		}
	}

	if endpoint.ID == 0 {
		database.DB.Create(&endpoint)
		util.Info("[endpoint] 添加端点: %s", endpoint.Tag)
	} else {
		var old database.Endpoint
		database.DB.First(&old, endpoint.ID)
		database.DB.Save(&endpoint)
		util.Info("[endpoint] 更新端点: %s", endpoint.Tag)

		if old.Tag != endpoint.Tag {
			route.CleanupRule("endpoint", old.Tag, endpoint.Tag)
		}
	}

	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"data": endpoint})
}

func DeleteEndpoint(c *gin.Context) {
	id := c.Param("id")

	var endpoint database.Endpoint
	database.DB.First(&endpoint, id)

	route.CleanupRule("endpoint", endpoint.Tag, "")

	if err := database.DB.Delete(&database.Endpoint{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	util.Info("[endpoint] 删除端点: %s", id)
	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ToggleEndpoint(c *gin.Context) {
	id := c.Param("id")
	var endpoint database.Endpoint
	if err := database.DB.First(&endpoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "端点不存在"})
		return
	}
	endpoint.Enable = !endpoint.Enable
	database.DB.Save(&endpoint)

	if !endpoint.Enable {
		route.CleanupRule("endpoint", endpoint.Tag, "")
	}

	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"data": endpoint})
}

func CreateWarpEndpoint(c *gin.Context) {
	var req struct {
		Address       string
		Reserved      string
		PrivateKey    string
		PublicKey     string
		PeerAddress   string
		PeerPort      int
		PeerPublicKey string
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var endpoint database.Endpoint
	database.DB.Where("tag = ?", "warp").First(&endpoint)

	if endpoint.ID == 0 {
		endpoint.Enable = true
		endpoint.Tag = "warp"
		endpoint.Type = "wireguard"
		endpoint.MTU = 1408
		endpoint.AllowedIPs = "0.0.0.0/0,::/0"
	}

	endpoint.Address = req.Address
	endpoint.Reserved = req.Reserved
	endpoint.PrivateKey = req.PrivateKey
	endpoint.PublicKey = req.PublicKey
	endpoint.PeerAddress = req.PeerAddress
	endpoint.PeerPort = req.PeerPort
	endpoint.PeerPublicKey = req.PeerPublicKey

	if err := database.DB.Save(&endpoint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	util.Info("[endpoint] 创建/更新 WARP 端点")
	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"data": endpoint})
}
