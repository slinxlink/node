package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	statsService "github.com/sagernet/sing-box/experimental/v2rayapi"
	"github.com/slinxlink/node/internal/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type Traffic struct {
	Upload   int64
	Download int64
}

var (
	boardMu      sync.Mutex
	boardTraffic = map[int]*Traffic{}
)

func queryTraffic() (users, inbounds map[string]*Traffic, err error) {
	conn, err := grpc.Dial("127.0.0.1:2048", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	req := &statsService.QueryStatsRequest{Reset_: true}
	resp := &statsService.QueryStatsResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.Invoke(ctx, "/v2ray.core.app.stats.command.StatsService/QueryStats", req, resp)
	if err != nil {
		return nil, nil, err
	}

	users = map[string]*Traffic{}
	inbounds = map[string]*Traffic{}

	for _, stat := range resp.Stat {
		parts := strings.Split(stat.Name, ">>>")
		if len(parts) != 4 {
			continue
		}

		var group map[string]*Traffic
		switch parts[0] {
		case "user":
			group = users
		case "inbound":
			group = inbounds
		default:
			continue
		}

		name := parts[1]
		if group[name] == nil {
			group[name] = &Traffic{}
		}

		switch parts[3] {
		case "uplink":
			group[name].Upload = stat.Value
		case "downlink":
			group[name].Download = stat.Value
		}
	}

	return users, inbounds, nil
}

func CollectTraffic() error {
	users, inbounds, err := queryTraffic()
	if err != nil {
		return err
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		for name, t := range users {
			if t.Upload == 0 && t.Download == 0 {
				continue
			}

			if uid := boardUserID(name); uid > 0 {
				addBoardTraffic(uid, t)
				continue
			}

			if err := addTraffic(tx.Model(&database.User{}).Where("name = ?", name), t); err != nil {
				return err
			}
		}

		for tag, t := range inbounds {
			port, err := strconv.Atoi(tag)
			if err != nil || (t.Upload == 0 && t.Download == 0) {
				continue
			}

			if err := addTraffic(tx.Model(&database.Inbound{}).Where("port = ?", port), t); err != nil {
				return err
			}
		}

		return nil
	})
}

func addTraffic(query *gorm.DB, t *Traffic) error {
	return query.UpdateColumns(map[string]interface{}{
		"upload":   gorm.Expr("COALESCE(upload, 0) + ?", t.Upload),
		"download": gorm.Expr("COALESCE(download, 0) + ?", t.Download),
	}).Error
}

func boardUserID(name string) int {
	var uid int
	if n, _ := fmt.Sscanf(name, "BoardUser_%d", &uid); n != 1 {
		return 0
	}
	return uid
}

func addBoardTraffic(uid int, t *Traffic) {
	boardMu.Lock()
	defer boardMu.Unlock()

	if boardTraffic[uid] == nil {
		boardTraffic[uid] = &Traffic{}
	}
	boardTraffic[uid].Upload += t.Upload
	boardTraffic[uid].Download += t.Download
}

func TakeBoardTraffic() map[int]Traffic {
	boardMu.Lock()
	defer boardMu.Unlock()

	result := make(map[int]Traffic, len(boardTraffic))
	for uid, t := range boardTraffic {
		result[uid] = *t
	}
	boardTraffic = map[int]*Traffic{}

	return result
}
