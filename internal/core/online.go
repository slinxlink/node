package core

import (
	"encoding/json"
	"sort"

	"github.com/slinxlink/node/internal/database"
	"gorm.io/gorm"
)

func CollectOnline() error {
	online, err := GetOnlineUsers()
	if err != nil {
		return err
	}

	ips := map[string][]string{}
	for _, u := range online {
		if boardUserID(u.Name) > 0 {
			continue
		}
		ips[u.Name] = append(ips[u.Name], u.IP)
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var users []database.User
		if err := tx.Find(&users).Error; err != nil {
			return err
		}

		for _, u := range users {
			want := "[]"

			if list := ips[u.Name]; len(list) > 0 {
				sort.Strings(list)
				data, _ := json.Marshal(list)
				want = string(data)
			}

			if u.AliveIP == want {
				continue
			}

			if err := tx.Model(&database.User{}).Where("id = ?", u.ID).UpdateColumn("alive_ip", want).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
