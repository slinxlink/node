package database

import (
	"os"
	"time"

	"github.com/seekky/slinx-node/internal/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() (bool, error) {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return false, err
	}

	DB, err = gorm.Open(sqlite.Open("data/slinx.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return false, err
	}

	err = DB.AutoMigrate(
		&Config{},
		&Core{},
		&Stats{},
		&Cert{},
		&Acme{},
		&DnsAccount{},
		&Inbound{},
		&User{},
		&Board{},
		&BoardUser{},
		&SystemLog{},
		&IP{},
		&Unlock{},
		&Route{},
	)
	if err != nil {
		return false, err
	}

	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_ip_source_version ON ip(source, ip_version)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unlock_ip_version_platform ON unlock(ip, ip_version, platform)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_route_city ON route(city)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_board_user ON board_user(board_id, user_id)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_core_single ON core(name)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_cert_domain ON cert(domain)")

	var coreCount int64
	DB.Model(&Core{}).Count(&coreCount)
	if coreCount == 0 {
		DB.Create(&Core{
			Name:       "sing-box",
			BinPath:    "bin/sing-box",
			ConfigPath: "data/sing-box.json",
			Repo:       "https://github.com/SagerNet/sing-box",
			LogEnable:  true,
			LogLevel:   "info",
			LogPath:    "data/sing-box.log",
		})
	}

	isFirstRun, err := initConfig()
	if err != nil {
		return false, err
	}

	initSystemLog()
	initStats()

	return isFirstRun, nil
}

func initConfig() (bool, error) {
	var count int64
	DB.Model(&Config{}).Count(&count)
	if count > 0 {
		return false, nil
	}

	ipv4, ipv6 := util.GetPublicIPs()

	config := Config{
		SecretKey:   util.GenerateString(32),
		Username:    "admin",
		Password:    util.GenerateString(12),
		Port:        util.GeneratePort(),
		Path:        "/" + util.GenerateString(8),
		IPv4:        ipv4,
		IPv6:        ipv6,
		LogEnable:   true,
		LogLevel:    "info",
		LogPath:     "data/slinx.log",
		BoardEnable: false,
	}

	return true, DB.Create(&config).Error
}

func initSystemLog() {
	var count int64
	DB.Model(&SystemLog{}).Count(&count)
	if count > 0 {
		return
	}

	now := time.Now()
	mins := now.Minute() / 15 * 15
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), mins, 0, 0, now.Location())
	base := now.Add(-24 * time.Hour)

	logs := make([]SystemLog, 96)
	for i := 0; i < 96; i++ {
		logs[i] = SystemLog{
			CreatedAt: base.Add(time.Duration(i) * 15 * time.Minute),
		}
	}
	DB.Create(&logs)
}

func initStats() {
	var count int64
	DB.Model(&Stats{}).Count(&count)
	if count > 0 {
		return
	}
	// 创建面板总流量记录
	DB.Create(&Stats{InboundID: 0})
}
