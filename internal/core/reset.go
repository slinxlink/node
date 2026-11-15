package core

import (
	"time"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

func ValidateReset(cycle string, day int) string {
	switch cycle {
	case "", "day":
		return ""
	case "week":
		if day < 1 || day > 7 {
			return "每周重置日必须在 1 到 7 之间"
		}
	case "month":
		if day < 1 || day > 31 {
			return "每月重置日必须在 1 到 31 之间"
		}
	default:
		return "重置周期无效"
	}
	return ""
}

func lastResetDue(cycle string, day int, now time.Time) time.Time {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch cycle {
	case "week":
		weekday := day % 7
		back := (int(today.Weekday()) - weekday + 7) % 7
		return today.AddDate(0, 0, -back)
	case "month":
		due := monthDue(now.Year(), now.Month(), day, now.Location())
		if due.After(now) {
			previous := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
			due = monthDue(previous.Year(), previous.Month(), day, now.Location())
		}
		return due
	}

	return today
}

func monthDue(year int, month time.Month, day int, location *time.Location) time.Time {
	last := time.Date(year, month+1, 0, 0, 0, 0, 0, location).Day()
	return time.Date(year, month, min(day, last), 0, 0, 0, 0, location)
}

func resetDue(cycle string, day int, resetAt, now time.Time) bool {
	return cycle != "" && resetAt.Before(lastResetDue(cycle, day, now))
}

func ResetTraffic() {
	now := time.Now()

	var inbounds []database.Inbound
	database.DB.Where("reset_cycle IS NOT NULL AND reset_cycle != ''").Find(&inbounds)
	for _, ib := range inbounds {
		if ib.ResetAt.IsZero() {
			database.DB.Model(&ib).Update("reset_at", now)
			continue
		}
		if resetDue(ib.ResetCycle, ib.ResetDay, ib.ResetAt, now) {
			database.DB.Model(&ib).Updates(map[string]any{"upload": 0, "download": 0, "reset_at": now})
			util.Info("[stats] 重置入站流量: %s:%d", ib.Protocol, ib.Port)
		}
	}

	var users []database.User
	database.DB.Where("reset_cycle IS NOT NULL AND reset_cycle != ''").Find(&users)
	for _, u := range users {
		if u.ResetAt.IsZero() {
			database.DB.Model(&u).Update("reset_at", now)
			continue
		}
		if resetDue(u.ResetCycle, u.ResetDay, u.ResetAt, now) {
			database.DB.Model(&u).Updates(map[string]any{"upload": 0, "download": 0, "reset_at": now})
			util.Info("[stats] 重置用户流量: %s", u.Name)
		}
	}
}
