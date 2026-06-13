package bootstrap

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seekky/slinx-node/internal/api"
	"github.com/seekky/slinx-node/internal/cli"
	"github.com/seekky/slinx-node/internal/config"
	"github.com/seekky/slinx-node/internal/core"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/job"
	syncer "github.com/seekky/slinx-node/internal/sync"
	"github.com/seekky/slinx-node/internal/util"
)

var WebFS embed.FS

func Start(webFS embed.FS) {
	WebFS = webFS

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

	// ── 5. 启动 Web 服务 ─────────────────────────────────────────────────────
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	api.RegisterRoutes(r)

	dist, _ := fs.Sub(webFS, "web/dist")
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/assets/") || path == "/favicon.ico" {
			http.FileServer(http.FS(dist)).ServeHTTP(c.Writer, c.Request)
			return
		}

		if path == config.Config.Path || strings.HasPrefix(path, config.Config.Path+"/") {
			data, err := fs.ReadFile(webFS, "web/dist/index.html")
			if err != nil {
				c.Status(500)
				return
			}
			html := strings.Replace(
				string(data),
				"</head>",
				fmt.Sprintf(`<script>window.__PANEL_PATH__ = '%s'</script></head>`, config.Config.Path),
				1,
			)
			c.Data(200, "text/html; charset=utf-8", []byte(html))
			return
		}

		c.Status(404)
	})

	if cfg.Domain != "" {
		var cert database.Cert
		database.DB.Where("domain = ?", cfg.Domain).First(&cert)
		go r.RunTLS(fmt.Sprintf(":%d", cfg.Port), cert.CertPath, cert.KeyPath)
	} else {
		go r.Run(fmt.Sprintf(":%d", cfg.Port))
	}

	// ── 6. 启动 sing-box 核心 ────────────────────────────────────────────────
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

	// ── 7. 启动面板对接同步 ──────────────────────────────────────────────────
	if cfg.BoardEnable {
		syncer.Start()
	}

	// ── 8. 启动后台任务 ──────────────────────────────────────────────────────
	job.Stats()
	job.SystemLog()
	job.CoreLogRotate()
	job.CertRenew()

	// ── 9. 收尾：更新启动时间、获取公网 IP ──────────────────────────────────
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
