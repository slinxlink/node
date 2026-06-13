package core

import (
	"context"
	"strings"

	statsService "github.com/sagernet/sing-box/experimental/v2rayapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserTraffic struct {
	Name     string
	Upload   int64
	Download int64
}

func GetUserTraffic() ([]UserTraffic, error) {
	conn, err := grpc.Dial("127.0.0.1:2048", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &statsService.QueryStatsRequest{
		Pattern: "user>>>",
		Reset_:  true,
	}
	resp := &statsService.QueryStatsResponse{}

	err = conn.Invoke(context.Background(), "/v2ray.core.app.stats.command.StatsService/QueryStats", req, resp)
	if err != nil {
		return nil, err
	}

	trafficMap := make(map[string]*UserTraffic)
	for _, stat := range resp.Stat {
		parts := strings.Split(stat.Name, ">>>")
		if len(parts) != 4 {
			continue
		}
		name := parts[1]
		direction := parts[3]

		if _, ok := trafficMap[name]; !ok {
			trafficMap[name] = &UserTraffic{Name: name}
		}
		if direction == "uplink" {
			trafficMap[name].Upload = stat.Value
		} else if direction == "downlink" {
			trafficMap[name].Download = stat.Value
		}
	}

	var result []UserTraffic
	for _, t := range trafficMap {
		result = append(result, *t)
	}
	return result, nil
}
