package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/cert"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/task"
	"github.com/slinxlink/node/internal/util"
)

// ── 证书 ─────────────────────────────────────────────────────────────────────

func GetCert(c *gin.Context) {
	var list []database.Cert
	database.DB.Find(&list)
	c.JSON(http.StatusOK, list)
}

func SaveCert(c *gin.Context) {
	var p struct {
		database.Cert
		CertContent string `json:"CertContent"`
		KeyContent  string `json:"KeyContent"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if p.Mode == "manual" {
		if p.CertContent == "" || p.KeyContent == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "证书内容或密钥内容不能为空"})
			return
		}
		if err := cert.ApplyManual(&p.Cert, p.CertContent, p.KeyContent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p.Cert)
		return
	}
	p.Domain = strings.TrimSpace(p.Domain)
	if !util.ValidateDomain(p.Domain) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名格式不正确"})
		return
	}
	if p.Acme == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择 ACME 账号"})
		return
	}
	if p.Mode == "dns" && p.Dns == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择 DNS 账号"})
		return
	}
	database.DB.Save(&p.Cert)
	c.JSON(http.StatusOK, p.Cert)
}

func DeleteCert(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dbCert database.Cert
	if err := database.DB.First(&dbCert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "证书不存在"})
		return
	}
	if err := cert.RemoveCert(&dbCert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func GetCertContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dbCert database.Cert
	if err := database.DB.First(&dbCert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "证书配置不存在"})
		return
	}

	content, err := cert.ReadCertContent(&dbCert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func ApplyCert(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dbCert database.Cert
	if err := database.DB.First(&dbCert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "证书不存在"})
		return
	}
	t := task.New(fmt.Sprintf("cert-%s", dbCert.Domain))
	go cert.ApplyCert(&dbCert, t)
	c.JSON(http.StatusOK, gin.H{"task_id": t.ID})
}

// ── ACME 账号 ─────────────────────────────────────────────────────────────────

func GetAcme(c *gin.Context) {
	var list []database.Acme
	database.DB.Find(&list)
	c.JSON(http.StatusOK, list)
}

func SaveAcme(c *gin.Context) {
	var p database.Acme
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if p.ID != 0 {
		var old database.Acme
		database.DB.First(&old, p.ID)
		if old.Email != p.Email || old.Provider != p.Provider || old.EabKid != p.EabKid || old.EabHmac != p.EabHmac {
			p.PrivateKey = ""
		}
	}
	if p.ID == 0 {
		database.DB.Create(&p)
		util.Info("[acme] 添加账号: %s", p.Email)
	} else {
		database.DB.Save(&p)
		util.Info("[acme] 更新账号: %s", p.Email)
	}
	if p.PrivateKey == "" {
		t := task.New(fmt.Sprintf("acme-%d", p.ID))
		go cert.RegisterAcme(&p, t)
		c.JSON(http.StatusOK, gin.H{"task_id": t.ID, "acme": p})
		return
	}
	c.JSON(http.StatusOK, p)
}

func DeleteAcme(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&database.Acme{}, id)
	util.Info("[acme] 删除账号: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// ── DNS 账号 ──────────────────────────────────────────────────────────────────

func GetDnsAccount(c *gin.Context) {
	var list []database.DnsAccount
	database.DB.Find(&list)
	c.JSON(http.StatusOK, list)
}

func SaveDnsAccount(c *gin.Context) {
	var p database.DnsAccount
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	p.Name = strings.TrimSpace(p.Name)
	p.Key = strings.TrimSpace(p.Key)
	p.Secret = strings.TrimSpace(p.Secret)
	if p.Name == "" || p.Key == "" || p.Secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整信息"})
		return
	}
	if p.ID == 0 {
		database.DB.Create(&p)
		util.Info("[dns] 添加 DNS 账号: %s", p.Name)
	} else {
		database.DB.Save(&p)
		util.Info("[dns] 更新 DNS 账号: %s", p.Name)
	}
	c.JSON(http.StatusOK, p)
}

func DeleteDnsAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&database.DnsAccount{}, id)
	util.Info("[dns] 删除 DNS 账号: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
