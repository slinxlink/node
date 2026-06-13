package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

func Quick(c *gin.Context) {
	usedPorts := database.UsedPorts()

	// ── 生成端口 ─────────────────────────────────────────────────────────
	var port int
	for {
		port = util.GeneratePort()
		if util.ValidatePort(port, usedPorts) == "" {
			break
		}
	}

	// ── 生成 Reality 密钥对 ───────────────────────────────────────────────
	privateKey, publicKey, err := util.GenerateRealityKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密钥对生成失败: " + err.Error()})
		return
	}

	// ── 生成 ShortIDs ─────────────────────────────────────────────────────
	shortIDs, err := util.GenerateShortIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ShortID 生成失败: " + err.Error()})
		return
	}
	shortIDsJSON, _ := json.Marshal(shortIDs)

	// ── 生成 UUID ─────────────────────────────────────────────────────────
	uuid, err := util.GenerateUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID 生成失败: " + err.Error()})
		return
	}

	// ── 创建入站 ─────────────────────────────────────────────────────────
	target := util.GenerateRealityTarget()
	inbound := database.Inbound{
		Enable:            true,
		Name:              "VLESS-" + time.Now().Format("0102-1504"),
		Protocol:          "vless",
		Port:              port,
		Transport:         "raw",
		TLSType:           "reality",
		RealityServerName: target,
		RealityServer:     target,
		RealityServerPort: 443,
		UTLS:              "chrome",
		RealityShortIDs:   string(shortIDsJSON),
		RealityPublicKey:  publicKey,
		RealityPrivateKey: privateKey,
		Flow:              "xtls-rprx-vision",
	}
	if err := database.DB.Create(&inbound).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "入站创建失败: " + err.Error()})
		return
	}

	// ── 创建用户 ─────────────────────────────────────────────────────────
	inboundIDs, _ := json.Marshal([]uint{inbound.ID})
	user := database.User{
		Name:     "用户-" + time.Now().Format("0102-1504"),
		Token:    util.GenerateString(32),
		Inbounds: string(inboundIDs),
		UUID:     uuid,
		Enable:   true,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败: " + err.Error()})
		return
	}

	// ── 应用配置 ─────────────────────────────────────────────────────────
	go core.Default.Apply()

	c.JSON(http.StatusOK, gin.H{
		"inbound": inbound,
		"user":    user,
	})
}
