package server

import (
	"context"
	"embed"
	"fmt"
	iofs "io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

// MARK: - 全局资源

var fs embed.FS

func Init(site embed.FS) {
	fs = site
}

var (
	Web = New("web")
	Sub = New("sub")
)

// MARK: - 服务实例

type Instance struct {
	name    string
	srv     *http.Server
	handler http.Handler
}

func New(name string) *Instance {
	return &Instance{name: name}
}

// MARK: - 服务启动

func (i *Instance) Start(handler http.Handler, config database.Config) error {
	i.handler = handler

	engine := handler.(*gin.Engine)
	dist, _ := iofs.Sub(fs, "web/dist")

	var port int

	switch i.name {
	case "web":
		port = config.Port

		engine.NoRoute(func(c *gin.Context) {
			if serveStatic(c, dist) {
				return
			}

			path := c.Request.URL.Path
			if path == config.Path || strings.HasPrefix(path, config.Path+"/") {
				data, err := iofs.ReadFile(fs, "web/dist/index.html")
				if err != nil {
					c.Status(500)
					return
				}

				html := strings.Replace(
					string(data),
					"</head>",
					fmt.Sprintf(`<script>window.__PANEL_PATH__ = '%s'</script></head>`, config.Path),
					1,
				)

				c.Data(200, "text/html; charset=utf-8", []byte(html))
				return
			}

			c.Status(404)
		})

	case "sub":
		if !config.SubEnable {
			return nil
		}

		port = config.SubPort

		engine.NoRoute(func(c *gin.Context) {
			serveStatic(c, dist)
			c.Status(404)
		})
	}

	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: addr, Handler: handler}

	certPath, keyPath, err := resolveCert(config.Domain)
	if err != nil {
		return err
	}

	i.srv = srv

	go func() {
		var err error

		if certPath != "" {
			err = srv.ListenAndServeTLS(certPath, keyPath)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			util.Error("[server:%s] 服务启动失败 (%s): %v", i.name, addr, err)
		}
	}()

	return nil
}

// MARK: - 服务停止 / 重启

func (i *Instance) Stop() error {
	if i.srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := i.srv.Shutdown(ctx); err != nil {
		util.Error("[server:%s] 关闭失败: %v", i.name, err)
		return err
	}

	i.srv = nil
	return nil
}

func (i *Instance) Restart(config database.Config) error {
	if err := i.Stop(); err != nil {
		return err
	}

	return i.Start(i.handler, config)
}

// MARK: - 证书更新

func UpdateCert(domain string) {
	var config database.Config
	database.DB.First(&config)

	if config.Domain != domain {
		return
	}

	Web.Restart(config)
	Sub.Restart(config)
}

// MARK: - 静态资源

func serveStatic(c *gin.Context, dist iofs.FS) bool {
	path := c.Request.URL.Path

	if strings.HasPrefix(path, "/assets/") || path == "/favicon.ico" {
		http.FileServer(http.FS(dist)).ServeHTTP(c.Writer, c.Request)
		return true
	}

	return false
}

// MARK: - 证书

func resolveCert(domain string) (certPath, keyPath string, err error) {
	if domain == "" {
		return "", "", nil
	}

	var cert database.Cert
	database.DB.Where("domain = ?", domain).First(&cert)

	if cert.CertPath == "" || cert.KeyPath == "" {
		return "", "", fmt.Errorf("域名 %s 证书路径未配置", domain)
	}

	return cert.CertPath, cert.KeyPath, nil
}

// MARK: - 引擎构筑

func BuildEngine(name string) *gin.Engine {
	var engine *gin.Engine

	switch name {
	case "web":
		engine = gin.Default()
	case "sub":
		engine = gin.New()
	default:
		return nil
	}

	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false

	return engine
}
