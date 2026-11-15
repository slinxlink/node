package core

import (
	"sync"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

type usableSet struct {
	inbounds map[uint]bool
	users    map[uint]bool
}

var (
	usableMu   sync.Mutex
	lastUsable *usableSet
)

func usableInbounds(list []database.Inbound) []database.Inbound {
	var result []database.Inbound
	for _, ib := range list {
		if ib.Usable() {
			result = append(result, ib)
		}
	}
	return result
}

func usableUsers(list []database.User) []database.User {
	var result []database.User
	for _, u := range list {
		if u.Usable() {
			result = append(result, u)
		}
	}
	return result
}

func snapshot(inbounds []database.Inbound, users []database.User) *usableSet {
	set := &usableSet{inbounds: map[uint]bool{}, users: map[uint]bool{}}
	for _, ib := range inbounds {
		set.inbounds[ib.ID] = true
	}
	for _, u := range users {
		set.users[u.ID] = true
	}
	return set
}

func (s *usableSet) equal(other *usableSet) bool {
	if len(s.inbounds) != len(other.inbounds) || len(s.users) != len(other.users) {
		return false
	}
	for id := range s.inbounds {
		if !other.inbounds[id] {
			return false
		}
	}
	for id := range s.users {
		if !other.users[id] {
			return false
		}
	}
	return true
}

func recordUsable(inbounds []database.Inbound, users []database.User) {
	usableMu.Lock()
	defer usableMu.Unlock()

	lastUsable = snapshot(inbounds, users)
}

func limitsChanged() bool {
	usableMu.Lock()
	last := lastUsable
	usableMu.Unlock()

	if last == nil {
		return false
	}

	var inbounds []database.Inbound
	var users []database.User
	database.DB.Where("enable = ?", true).Find(&inbounds)
	database.DB.Where("enable = ?", true).Find(&users)

	return !snapshot(usableInbounds(inbounds), usableUsers(users)).equal(last)
}

func CheckLimits() {
	if Default.Status() == "stopped" || !limitsChanged() {
		return
	}

	util.Info("[core] 有入站或用户达到流量上限或有效期，重新生成配置并重启核心")

	if err := Default.Apply(); err != nil {
		util.Error("[core] 应用配置失败: %v", err)
	}
}
