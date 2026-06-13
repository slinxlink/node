package slinx

import "encoding/json"

type Response struct {
	Ret  uint            `json:"ret"`
	Data json.RawMessage `json:"data"`
}

type UserResponse struct {
	ID      int    `json:"id"`
	Passwd  string `json:"passwd"`
	UUID    string `json:"uuid"`
	AliveIP int    `json:"alive_ip"`
}

type PostData struct {
	Data interface{} `json:"data"`
}

type OnlineUser struct {
	UID int    `json:"user_id"`
	IP  string `json:"ip"`
}

type UserTraffic struct {
	UID      int   `json:"user_id"`
	Upload   int64 `json:"u"`
	Download int64 `json:"d"`
}
