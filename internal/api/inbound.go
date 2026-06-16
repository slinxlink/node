package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

func GetInbounds(c *gin.Context) {
	var inbounds []database.Inbound
	database.DB.Find(&inbounds)
	c.JSON(http.StatusOK, inbounds)
}

func SaveInbound(c *gin.Context) {
	var ib database.Inbound
	if err := c.ShouldBindJSON(&ib); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usedPorts := database.UsedPorts(ib.ID)

	if msg := util.ValidatePort(ib.Port, usedPorts); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	if ib.ObfsType != "" {
		if ib.ObfsPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写混淆密码"})
			return
		}
		if ib.ObfsType == "gecko" {
			if ib.ObfsMinPacketSize <= 0 || ib.ObfsMaxPacketSize <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "包大小必须大于 0"})
				return
			}
			if ib.ObfsMaxPacketSize <= ib.ObfsMinPacketSize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "最大包大小必须大于最小包大小"})
				return
			}
		}
	}

	if ib.TLSType == "tls" {
		var ids []int
		json.Unmarshal([]byte(ib.Certs), &ids)
		if len(ids) == 0 || ids[0] == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择证书"})
			return
		}
	}

	if ib.TLSType == "reality" {
		if ib.RealityPrivateKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写私钥"})
			return
		}
		if ib.RealityShortIDs == "" || ib.RealityShortIDs == "[]" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写短 ID"})
			return
		}
		if ib.RealityServer == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写目标服务器"})
			return
		}
		if ib.RealityServerPort == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写目标端口"})
			return
		}
		if ib.RealityServerName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写 SNI"})
			return
		}
	}

	if ib.ID == 0 {
		database.DB.Create(&ib)
		util.Info("[inbound] 添加入站: %s:%d", ib.Protocol, ib.Port)
	} else {
		database.DB.Save(&ib)
		util.Info("[inbound] 更新入站: %s:%d", ib.Protocol, ib.Port)
	}

	go core.Default.Apply()
	c.JSON(http.StatusOK, ib)
}

func DeleteInbound(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 清理 User 表里关联该入站的引用
	var users []database.User
	database.DB.Find(&users)
	for _, u := range users {
		var ids []int
		json.Unmarshal([]byte(u.Inbounds), &ids)
		newIDs := []int{}
		for _, i := range ids {
			if i != id {
				newIDs = append(newIDs, i)
			}
		}
		updated, _ := json.Marshal(newIDs)
		database.DB.Model(&u).Update("inbounds", string(updated))
	}

	// 清理 Board 表里关联该入站的引用
	database.DB.Where("inbound = ?", id).Update("inbound", 0)

	database.DB.Delete(&database.Inbound{}, id)
	util.Info("[inbound] 删除入站: %d", id)
	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ToggleInbound(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var ib database.Inbound
	if database.DB.First(&ib, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存在"})
		return
	}
	ib.Enable = !ib.Enable
	database.DB.Save(&ib)
	util.Info("[inbound] %s入站: %s:%d", map[bool]string{true: "启用", false: "禁用"}[ib.Enable], ib.Protocol, ib.Port)
	go core.Default.Apply()
	c.JSON(http.StatusOK, ib)
}

func QuickInbound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
