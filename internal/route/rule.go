package route

import (
	"encoding/json"

	"github.com/slinxlink/node/internal/database"
)

// CleanupRule 统一入口，根据 type 分流
func CleanupRule(typ, tag, newTag string) {
	switch typ {
	case "inbound":
		if newTag != "" {
			renameInbound(tag, newTag)
		} else {
			deleteInbound(tag)
		}
	case "endpoint":
		if newTag != "" {
			renameEndpoint(tag, newTag)
		} else {
			deleteEndpoint(tag)
		}
	}
}

// deleteSortGroup 删掉指定 sort 下的所有规则行（inbound/endpoint 共用）
func deleteSortGroup(sort int) {
	database.DB.Where("sort = ?", sort).Delete(&database.Rule{})
}

// renameInbound 把 inbound 数组里的旧 tag 替换成新 tag
func renameInbound(tag, newTag string) {
	var rows []database.Rule
	database.DB.Where("`key` = ?", "inbound").Find(&rows)

	for _, row := range rows {
		var tags []string
		if err := json.Unmarshal([]byte(row.Value), &tags); err != nil {
			continue
		}
		changed := false
		for i, t := range tags {
			if t == tag {
				tags[i] = newTag
				changed = true
			}
		}
		if changed {
			data, _ := json.Marshal(tags)
			database.DB.Model(&database.Rule{}).Where("id = ?", row.ID).Update("value", string(data))
		}
	}
}

// deleteInbound 从 inbound 数组里删掉这个 tag，数组空了就删整组
func deleteInbound(tag string) {
	var rows []database.Rule
	database.DB.Where("`key` = ?", "inbound").Find(&rows)

	for _, row := range rows {
		var tags []string
		if err := json.Unmarshal([]byte(row.Value), &tags); err != nil {
			continue
		}
		var remain []string
		for _, t := range tags {
			if t != tag {
				remain = append(remain, t)
			}
		}
		if len(remain) == len(tags) {
			continue // 没有这个 tag，跳过
		}
		if len(remain) == 0 {
			deleteSortGroup(row.Sort)
		} else {
			data, _ := json.Marshal(remain)
			database.DB.Model(&database.Rule{}).Where("id = ?", row.ID).Update("value", string(data))
		}
	}
}

// renameEndpoint 把 outbound 指向旧 tag 的那一行改成新 tag
func renameEndpoint(tag, newTag string) {
	tagJSON, _ := json.Marshal(tag)
	newTagJSON, _ := json.Marshal(newTag)
	database.DB.Model(&database.Rule{}).
		Where("`key` = ? AND value = ?", "outbound", string(tagJSON)).
		Update("value", string(newTagJSON))
}

// deleteEndpoint 找到 outbound 指向这个 tag 的所有 sort，整组删掉
func deleteEndpoint(tag string) {
	tagJSON, _ := json.Marshal(tag)
	var rows []database.Rule
	database.DB.Where("`key` = ? AND value = ?", "outbound", string(tagJSON)).Find(&rows)
	for _, row := range rows {
		deleteSortGroup(row.Sort)
	}
}
