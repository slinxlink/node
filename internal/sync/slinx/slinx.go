package slinx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/seekky/slinx-node/internal/core"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/util"
	"gorm.io/gorm"
)

func Sync(b database.Board) {
	users, err := getUsers(b)
	if err != nil {
		util.Error("[board:%s] 拉取用户失败: %v", b.Name, err)
		return
	}

	var existing []database.BoardUser
	database.DB.Where("board_id = ?", b.ID).Find(&existing)

	incomingIDs := make(map[int]bool)
	for _, u := range users {
		incomingIDs[u.ID] = true
	}
	for _, e := range existing {
		if !incomingIDs[e.UserID] {
			database.DB.Delete(&e)
		}
	}

	changed := diff(existing, users)

	for _, u := range users {
		bu := database.BoardUser{
			BoardID:   b.ID,
			UserID:    u.ID,
			UUID:      u.UUID,
			Passwd:    u.Passwd,
			AliveIP:   u.AliveIP,
			UpdatedAt: time.Now(),
		}
		database.DB.Where(database.BoardUser{BoardID: b.ID, UserID: u.ID}).
			Assign(bu).
			FirstOrCreate(&bu)
	}

	if changed {
		util.Info("[board:%s] 用户列表有变化，重新生成配置", b.Name)
		core.Default.Apply()
	}

	if core.Default.Status() == "stopped" {
		util.Info("[board:%s] 拉取到 %d 个用户", b.Name, len(users))
		return
	}

	traffic, err := core.GetUserTraffic()
	if err != nil {
		util.Error("[board:%s] 获取流量失败: %v", b.Name, err)
		return
	}

	onlineUsers, err := core.GetOnlineUsers()
	if err != nil {
		util.Error("[board:%s] 获取在线用户失败: %v", b.Name, err)
	}

	util.Info("[board:%s] 拉取到 %d 个用户，上报 %d 条流量 %d 个在线IP", b.Name, len(users), len(traffic), len(onlineUsers))

	if len(traffic) > 0 {
		reportTraffic(b, traffic)
	}

	if err == nil {
		reportAliveIP(b, onlineUsers)
	}
}

func diff(existing []database.BoardUser, incoming []UserResponse) bool {
	if len(existing) != len(incoming) {
		return true
	}

	existMap := make(map[int]database.BoardUser)
	for _, u := range existing {
		existMap[u.UserID] = u
	}

	for _, u := range incoming {
		e, ok := existMap[u.ID]
		if !ok {
			return true
		}
		if u.UUID != "" && e.UUID != u.UUID {
			return true
		}
		if u.Passwd != "" && e.Passwd != u.Passwd {
			return true
		}
	}

	return false
}

func getUsers(b database.Board) ([]UserResponse, error) {
	url := fmt.Sprintf("%s/mod_mu/users?node_id=%d&key=%s", b.Host, b.NodeID, b.Key)

	res, err := http.Get(url)
	if err != nil {
		util.Error("[board:%s] 请求失败: %v", b.Name, err)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer res.Body.Close()

	var result Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		util.Error("[board:%s] 解析响应失败: %v", b.Name, err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Ret != 1 {
		util.Warn("[board:%s] 接口返回异常: ret=%d", b.Name, result.Ret)
		return nil, fmt.Errorf("接口返回异常: ret=%d", result.Ret)
	}

	var users []UserResponse
	if err := json.Unmarshal(result.Data, &users); err != nil {
		util.Error("[board:%s] 解析用户列表失败: %v", b.Name, err)
		return nil, fmt.Errorf("解析用户列表失败: %w", err)
	}

	return users, nil
}

func reportTraffic(b database.Board, traffic []core.UserTraffic) {
	trafficMap := make(map[int]core.UserTraffic)
	for _, t := range traffic {
		var uid int
		fmt.Sscanf(t.Name, "BoardUser_%d", &uid)
		if uid > 0 {
			trafficMap[uid] = t
		}
	}

	var data []UserTraffic
	for uid, t := range trafficMap {
		data = append(data, UserTraffic{
			UID:      uid,
			Upload:   t.Upload,
			Download: t.Download,
		})
		database.DB.Model(&database.BoardUser{}).
			Where("board_id = ? AND user_id = ?", b.ID, uid).
			Updates(map[string]interface{}{
				"upload":   gorm.Expr("upload + ?", t.Upload),
				"download": gorm.Expr("download + ?", t.Download),
			})
	}

	if len(data) == 0 {
		return
	}

	url := fmt.Sprintf("%s/mod_mu/users/traffic?node_id=%d&key=%s", b.Host, b.NodeID, b.Key)
	body, _ := json.Marshal(PostData{Data: data})
	http.Post(url, "application/json", bytes.NewReader(body))
}

func reportAliveIP(b database.Board, users []core.OnlineUser) {
	database.DB.Model(&database.BoardUser{}).
		Where("board_id = ?", b.ID).
		Update("alive_ip", 0)

	var data []OnlineUser
	for _, u := range users {
		var uid int
		fmt.Sscanf(u.Name, "BoardUser_%d", &uid)
		if uid > 0 {
			data = append(data, OnlineUser{
				UID: uid,
				IP:  u.IP,
			})
			database.DB.Model(&database.BoardUser{}).
				Where("board_id = ? AND user_id = ?", b.ID, uid).
				Update("alive_ip", u.IP)
		}
	}

	if len(data) == 0 {
		return
	}

	url := fmt.Sprintf("%s/mod_mu/users/aliveip?node_id=%d&key=%s", b.Host, b.NodeID, b.Key)
	body, _ := json.Marshal(PostData{Data: data})
	http.Post(url, "application/json", bytes.NewReader(body))
}
