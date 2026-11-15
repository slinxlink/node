package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

func GetUsers(c *gin.Context) {
	var users []database.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func SaveUser(c *gin.Context) {
	var u database.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(u.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	if u.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token 不能为空"})
		return
	}

	if u.TrafficLimit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "流量上限不能小于 0"})
		return
	}

	if msg := core.ValidateReset(u.ResetCycle, u.ResetDay); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var count int64
	database.DB.Model(&database.User{}).Where("token = ? AND id != ?", u.Token, u.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token 已存在"})
		return
	}

	database.DB.Model(&database.User{}).Where("name = ? AND id != ?", u.Name, u.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	if u.ID == 0 {
		u.ResetAt = time.Now()
		database.DB.Create(&u)
		util.Info("[user] 添加用户: %s", u.Name)
	} else {
		var old database.User
		database.DB.First(&old, u.ID)
		database.DB.Omit("Upload", "Download", "AliveIP", "ResetAt").Save(&u)
		if old.ResetCycle != u.ResetCycle || old.ResetDay != u.ResetDay {
			database.DB.Model(&u).Update("reset_at", time.Now())
		}
		util.Info("[user] 更新用户: %s", u.Name)
	}

	go core.Default.Apply()
	c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&database.User{}, id)
	util.Info("[user] 删除用户: %d", id)
	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ToggleUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var u database.User
	if database.DB.First(&u, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存在"})
		return
	}
	u.Enable = !u.Enable
	database.DB.Model(&u).Update("enable", u.Enable)
	util.Info("[user] %s用户: %s", map[bool]string{true: "启用", false: "禁用"}[u.Enable], u.Name)
	go core.Default.Apply()
	c.JSON(http.StatusOK, u)
}

func ResetUserTraffic(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var u database.User
	if database.DB.First(&u, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存在"})
		return
	}
	database.DB.Model(&u).Updates(map[string]any{"upload": 0, "download": 0})
	util.Info("[user] 重置用户流量: %s", u.Name)
	go core.CheckLimits()
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
