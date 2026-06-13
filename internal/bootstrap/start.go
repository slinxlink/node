package bootstrap

import (
	"fmt"
	"log"
	"time"

	"github.com/slinxlink/node/internal/cli"
	"github.com/slinxlink/node/internal/config"
	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/job"
	"github.com/slinxlink/node/internal/server"
	syncer "github.com/slinxlink/node/internal/sync"
	"github.com/slinxlink/node/internal/util"
)

func Start() {
	// ── 1. 数据库 ────────────────────────────────────────────────────────────
	isFirstRun, err := database.Init()
	if err != nil {
		cli.Status("数据库", "初始化失败", false)
		log.Fatal(err)
	}

	// ── 2. 读取配置 ──────────────────────────────────────────────────────────
	var cfg database.Config
	database.DB.First(&cfg)

	if err := config.Load(); err != nil {
		cli.Status("配置", "加载失败", false)
		log.Fatal(err)
	}

	// ── 3. 初始化日志 ────────────────────────────────────────────────────────
	util.InitLog(cfg.LogPath, cfg.LogLevel, cfg.LogEnable)

	// ── 4. 初始化核心 ────────────────────────────────────────────────────────
	core.Default.Init()

	// ── 5. 启动Web服务 ─────────────────────────────────────────────────────
	if err := server.StartWeb(); err != nil {
		util.Error("[server] Web服务启动失败: %v", err)
	}

	// ── 6 启动订阅服务 ─────────────────────────────────────────────────────
	if err := server.StartSub(); err != nil {
		util.Error("[server] 订阅服务启动失败: %v", err)
	}

	// ── 7. 启动 sing-box 核心 ────────────────────────────────────────────────
	if err := core.Default.Start(); err != nil {
		cli.Status("核心", "启动失败", false)
		log.Println(err)
	} else {
		cli.Status("核心", "启动成功", true)
		database.DB.Model(&database.Core{}).Where("id = 1").Updates(map[string]interface{}{
			"version":    core.Default.Version(),
			"started_at": time.Now(),
		})
	}

	// ── 8. 启动面板对接同步 ──────────────────────────────────────────────────
	if cfg.BoardEnable {
		syncer.Start()
	}

	// ── 9. 启动后台任务 ──────────────────────────────────────────────────────
	job.Stats()
	job.SystemLog()
	job.CoreLogRotate()
	job.CertRenew()
	job.ClashTemplateRefresh()

	// ── 10. 收尾：更新启动时间、获取公网 IP ──────────────────────────────────
	database.DB.Model(&database.Config{}).Where("id = 1").Update("started_at", time.Now())

	go func() {
		ipv4, ipv6 := util.GetPublicIPs()
		database.DB.Model(&database.Config{}).Where("id = 1").Updates(map[string]interface{}{
			"ipv4": ipv4,
			"ipv6": ipv6,
		})
		if isFirstRun {
			printFirstRun(ipv4)
		}
	}()

	cli.Status("面板", "启动成功", true)
	util.Info("[bootstrap] 面板启动")

	Shutdown()
}

func printFirstRun(ipv4 string) {
	host := ipv4
	if host == "" {
		host = "localhost"
	}
	addr := fmt.Sprintf("http://%s:%d%s", host, config.Config.Port, config.Config.Path)
	cli.Info("登录信息",
		[]string{"访问地址", addr},
		[]string{"用户名", config.Config.Username},
		[]string{"密码", config.Config.Password},
	)
}
