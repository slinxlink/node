package slinx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
	"gorm.io/gorm"
)

var client = &http.Client{Timeout: 30 * time.Second}

var (
	pendingMu sync.Mutex
	pending   = map[uint]map[int]*ReportUser{}
)

func Sync(b database.Board) {
	report := collect(b)

	users, err := post(b, report)
	if err != nil {
		util.Error("[board:%s] 同步失败: %v", b.Name, err)
		keep(b.ID, report)
		return
	}

	apply(b, users)

	var ips int
	for _, r := range report {
		ips += len(r.IP)
	}
	util.Info("[board:%s] 拉取到 %d 个用户，上报 %d 个用户的流量和在线情况 %d 个在线IP", b.Name, len(users), len(report), ips)
}

func collect(b database.Board) []ReportUser {
	usage := take(b.ID)

	if core.Default.Status() != "stopped" {
		traffic, err := core.GetUserTraffic()
		if err != nil {
			util.Error("[board:%s] 获取流量失败: %v", b.Name, err)
		}

		for _, t := range traffic {
			var uid int
			fmt.Sscanf(t.Name, "BoardUser_%d", &uid)
			if uid <= 0 || (t.Upload == 0 && t.Download == 0) {
				continue
			}

			entry(usage, uid).Upload += t.Upload
			entry(usage, uid).Download += t.Download

			database.DB.Model(&database.BoardUser{}).
				Where("board_id = ? AND user_id = ?", b.ID, uid).
				Updates(map[string]interface{}{
					"upload":   gorm.Expr("upload + ?", t.Upload),
					"download": gorm.Expr("download + ?", t.Download),
				})
		}

		online, err := core.GetOnlineUsers()
		if err != nil {
			util.Error("[board:%s] 获取在线用户失败: %v", b.Name, err)
		} else {
			database.DB.Model(&database.BoardUser{}).Where("board_id = ?", b.ID).Update("alive_ip", "[]")

			ips := map[int][]string{}
			for _, u := range online {
				var uid int
				fmt.Sscanf(u.Name, "BoardUser_%d", &uid)
				if uid > 0 {
					ips[uid] = append(ips[uid], u.IP)
				}
			}

			for uid, list := range ips {
				entry(usage, uid).IP = list

				ipJSON, _ := json.Marshal(list)
				database.DB.Model(&database.BoardUser{}).
					Where("board_id = ? AND user_id = ?", b.ID, uid).
					Update("alive_ip", string(ipJSON))
			}
		}
	}

	report := []ReportUser{}
	for _, r := range usage {
		report = append(report, *r)
	}
	return report
}

func entry(usage map[int]*ReportUser, uid int) *ReportUser {
	if usage[uid] == nil {
		usage[uid] = &ReportUser{ID: uid, IP: []string{}}
	}
	return usage[uid]
}

func take(boardID uint) map[int]*ReportUser {
	pendingMu.Lock()
	defer pendingMu.Unlock()

	usage := pending[boardID]
	delete(pending, boardID)

	if usage == nil {
		usage = map[int]*ReportUser{}
	}
	return usage
}

func keep(boardID uint, report []ReportUser) {
	pendingMu.Lock()
	defer pendingMu.Unlock()

	usage := map[int]*ReportUser{}
	for _, r := range report {
		if r.Upload > 0 || r.Download > 0 {
			usage[r.ID] = &ReportUser{ID: r.ID, Upload: r.Upload, Download: r.Download, IP: []string{}}
		}
	}

	if len(usage) > 0 {
		pending[boardID] = usage
	}
}

func post(b database.Board, report []ReportUser) ([]UserResponse, error) {
	body, err := json.Marshal(SyncRequest{User: report})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(b.Host, "/")+"/api/node/sync", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Node-Id", strconv.Itoa(b.NodeID))
	req.Header.Set("Node-Key", b.Key)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var e ErrorResponse
		json.NewDecoder(res.Body).Decode(&e)
		return nil, fmt.Errorf("HTTP %d %s", res.StatusCode, e.Message)
	}

	var result SyncResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return result.User, nil
}

func apply(b database.Board, users []UserResponse) {
	var existing []database.BoardUser
	database.DB.Where("board_id = ?", b.ID).Find(&existing)

	incoming := make(map[int]bool)
	for _, u := range users {
		incoming[u.ID] = true
	}
	for _, e := range existing {
		if !incoming[e.UserID] {
			database.DB.Delete(&e)
		}
	}

	changed := diff(existing, users)

	for _, u := range users {
		var bu database.BoardUser
		database.DB.Where(database.BoardUser{BoardID: b.ID, UserID: u.ID}).
			Attrs(database.BoardUser{AliveIP: "[]"}).
			Assign(database.BoardUser{UUID: u.UUID, Passwd: u.Password, UpdatedAt: time.Now()}).
			FirstOrCreate(&bu)
	}

	if changed {
		util.Info("[board:%s] 用户列表有变化，重新生成配置", b.Name)
		core.Default.Apply()
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
		if u.Password != "" && e.Passwd != u.Password {
			return true
		}
	}

	return false
}
