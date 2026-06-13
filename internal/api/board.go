package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	syncer "github.com/slinxlink/node/internal/sync"
	"github.com/slinxlink/node/internal/util"
)

func GetBoards(c *gin.Context) {
	var boards []database.Board
	database.DB.Find(&boards)
	c.JSON(http.StatusOK, boards)
}

func SaveBoard(c *gin.Context) {
	var p database.Board
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if p.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "面板地址不能为空"})
		return
	}
	if p.NodeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "节点 ID 不能为 0"})
		return
	}
	if p.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "通讯密钥不能为空"})
		return
	}
	if p.SyncInterval < 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "同步间隔不能小于 30 秒"})
		return
	}

	if p.ID == 0 {
		database.DB.Create(&p)
		util.Info("[board] 添加面板: %s", p.Name)
	} else {
		database.DB.Save(&p)
		util.Info("[board] 更新面板: %s", p.Name)
	}

	syncer.Restart()
	go core.Default.Apply()
	c.JSON(http.StatusOK, p)
}

func DeleteBoard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Where("board_id = ?", id).Delete(&database.BoardUser{})
	database.DB.Delete(&database.Board{}, id)
	util.Info("[board] 删除面板: %d", id)
	syncer.Restart()
	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ToggleBoard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p database.Board
	if database.DB.First(&p, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存在"})
		return
	}
	p.Enable = !p.Enable
	if !p.Enable {
		database.DB.Where("board_id = ?", id).Delete(&database.BoardUser{})
	}
	database.DB.Save(&p)
	util.Info("[board] %s面板: %s", map[bool]string{true: "启用", false: "禁用"}[p.Enable], p.Name)
	syncer.Restart()
	go core.Default.Apply()
	c.JSON(http.StatusOK, p)
}

func GetBoardUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var users []database.BoardUser
	database.DB.Where("board_id = ?", id).Find(&users)
	c.JSON(http.StatusOK, users)
}
