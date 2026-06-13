package api

import (
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seekky/slinx-node/internal/config"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/util"
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Username != config.Config.Username || req.Password != config.Config.Password {
		util.Warn("[auth] 登录失败，IP: %s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(config.Config.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
		return
	}

	util.Info("[auth] 登录成功，IP: %s", c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func Logout(c *gin.Context) {
	util.Info("[auth] 已登出")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ChangeCredentials(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewUsername string `json:"new_username"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.NewUsername = strings.TrimSpace(req.NewUsername)

	if req.OldPassword != config.Config.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "原密码错误"})
		return
	}

	if len(req.NewPassword) > 0 && len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能少于6位"})
		return
	}

	if len(req.NewUsername) > 0 {
		if len(req.NewUsername) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "新用户名不能少于6位"})
			return
		}
		for _, r := range req.NewUsername {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				c.JSON(http.StatusBadRequest, gin.H{"error": "用户名只能包含字母、数字和下划线"})
				return
			}
		}
	}

	updates := map[string]any{}
	if req.NewUsername != "" {
		updates["username"] = req.NewUsername
		config.Config.Username = req.NewUsername
	}
	if req.NewPassword != "" {
		updates["password"] = req.NewPassword
		config.Config.Password = req.NewPassword
	}

	database.DB.Model(&config.Config).Updates(updates)
	util.Info("[auth] 管理员凭据已更新")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
