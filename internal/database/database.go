package database

import (
	"os"
	"time"

	"github.com/slinxlink/node/internal/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() (bool, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return false, err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open("data/slinx.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return false, err
	}

	if err = DB.AutoMigrate(
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
		&Endpoint{},
		&Rule{},
		&Warp{},
		&SystemLog{},
		&IP{},
		&Unlock{},
		&BackRoute{},
	); err != nil {
		return false, err
	}

	createIndexes()
	initCore()
	initSystemLog()
	initStats()

	isFirstRun, err := initConfig()
	if err != nil {
		return false, err
	}

	if !isFirstRun {
		patchDefaults()
	}

	return isFirstRun, nil
}

func createIndexes() {
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_board_user ON board_user(board_id, user_id)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_endpoint_tag ON endpoint(tag)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_ip_source_version ON ip(source, ip_version)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unlock_ip_version_platform ON unlock(ip, ip_version, platform)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_back_route_city ON back_route(city)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_core_single ON core(name)")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_cert_domain ON cert(domain)")
}

func initConfig() (bool, error) {
	var count int64
	DB.Model(&Config{}).Count(&count)
	if count > 0 {
		return false, nil
	}

	ipv4, ipv6 := util.GetPublicIPs()

	config := Config{
		Username:  "admin",
		Password:  util.GenerateString(12),
		SecretKey: util.GenerateString(32),

		Port: util.GeneratePort(),
		Path: "/" + util.GenerateString(8),
		IPv4: ipv4,
		IPv6: ipv6,

		SubEnable:         true,
		SubPath:           "/link",
		SubPort:           2096,
		RulesetAutoUpdate: true,

		LogEnable: true,
		LogLevel:  "info",
		LogPath:   "data/slinx.log",

		BBR: true,

		BoardEnable: false,

		Repo: "https://github.com/slinxlink/node",
	}

	return true, DB.Create(&config).Error
}

func initCore() {
	var count int64
	DB.Model(&Core{}).Count(&count)
	if count > 0 {
		return
	}
	DB.Create(&Core{
		Name:       "sing-box",
		BinPath:    "bin/sing-box",
		ConfigPath: "data/sing-box.json",
		LogEnable:  true,
		LogLevel:   "info",
		LogPath:    "data/sing-box.log",
	})
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
	DB.Model(&Stats{}).Where("inbound_id = 0").Count(&count)
	if count > 0 {
		return
	}
	DB.Create(&Stats{InboundID: 0})
}

func patchDefaults() {
	patchStr := func(field *string, val string, dirty *bool) {
		if *field == "" {
			*field = val
			*dirty = true
		}
	}
	patchInt := func(field *int, val int, dirty *bool) {
		if *field == 0 {
			*field = val
			*dirty = true
		}
	}

	var cfg Config
	DB.First(&cfg)
	var cfgDirty bool

	patchStr(&cfg.Username, "admin", &cfgDirty)
	patchStr(&cfg.Password, util.GenerateString(12), &cfgDirty)
	patchStr(&cfg.SecretKey, util.GenerateString(32), &cfgDirty)
	patchInt(&cfg.Port, util.GeneratePort(), &cfgDirty)
	patchStr(&cfg.Path, "/"+util.GenerateString(8), &cfgDirty)
	patchStr(&cfg.SubPath, "/link", &cfgDirty)
	patchInt(&cfg.SubPort, 2096, &cfgDirty)
	patchStr(&cfg.LogLevel, "info", &cfgDirty)
	patchStr(&cfg.LogPath, "data/slinx.log", &cfgDirty)
	patchStr(&cfg.Repo, "https://github.com/slinxlink/node", &cfgDirty)
	if cfgDirty {
		DB.Save(&cfg)
	}

	var core Core
	DB.First(&core)
	var coreDirty bool

	patchStr(&core.BinPath, "bin/sing-box", &coreDirty)
	patchStr(&core.ConfigPath, "data/sing-box.json", &coreDirty)
	patchStr(&core.LogLevel, "info", &coreDirty)
	patchStr(&core.LogPath, "data/sing-box.log", &coreDirty)
	if coreDirty {
		DB.Save(&core)
	}
}
