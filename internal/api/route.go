package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

// GET /rule 整表返回
func GetRule(c *gin.Context) {
	var rules []database.Rule
	database.DB.Order("sort asc, `index` asc").Find(&rules)
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func SaveRule(c *gin.Context) {
	var req struct {
		Tag      string   `json:"tag"`
		Inbounds []string `json:"inbounds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	tagJSON, _ := json.Marshal(req.Tag)
	inboundJSON, _ := json.Marshal(req.Inbounds)

	var outboundRow database.Rule
	err := database.DB.Where("`key` = ? AND value = ?", "outbound", string(tagJSON)).First(&outboundRow).Error

	if err == nil {
		// 已存在对应 tag 的组
		if len(req.Inbounds) == 0 {
			// 清空了，直接删掉整组
			database.DB.Where("sort = ?", outboundRow.Sort).Delete(&database.Rule{})
		} else {
			database.DB.Model(&database.Rule{}).
				Where("sort = ? AND `key` = ?", outboundRow.Sort, "inbound").
				Update("value", string(inboundJSON))
		}
	} else if len(req.Inbounds) > 0 {
		// 不存在,且这次有值才新开一组（空列表不需要新建）
		var maxSort int
		database.DB.Model(&database.Rule{}).Select("COALESCE(MAX(sort),0)").Scan(&maxSort)
		newSort := maxSort + 1

		database.DB.Create(&database.Rule{Sort: newSort, Index: 0, Key: "inbound", Value: string(inboundJSON)})
		database.DB.Create(&database.Rule{Sort: newSort, Index: 1, Key: "outbound", Value: string(tagJSON)})
	}

	util.Info("[rule] 更新路由: %s", req.Tag)
	go core.Default.Apply()
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
