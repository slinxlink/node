package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slinxlink/node/internal/api"
	"github.com/slinxlink/node/internal/database"
)

var webFS embed.FS

func Init(fs embed.FS) {
	webFS = fs
}

func StartWeb() error {
	var config database.Config
	database.DB.First(&config)

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

		if path == config.Path || strings.HasPrefix(path, config.Path+"/") {
			data, err := fs.ReadFile(webFS, "web/dist/index.html")
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

	return runEngine(r, config.Domain, config.Port)
}

func StartSub() error {
	var config database.Config
	database.DB.First(&config)

	if !config.SubEnable {
		return nil
	}

	subEngine := gin.New()
	subDist, _ := fs.Sub(webFS, "web/dist")
	subEngine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/assets/") || path == "/favicon.ico" {
			http.FileServer(http.FS(subDist)).ServeHTTP(c.Writer, c.Request)
		}
	})

	subEngine.GET(config.SubPath+"/:token", func(c *gin.Context) {
		accept := c.GetHeader("Accept")
		if strings.Contains(strings.ToLower(accept), "text/html") {
			api.GetSubscriptionPage(c, webFS)
		} else {
			api.GetSubscription(c)
		}
	})

	subEngine.GET(config.SubPath+"/:token/clash", func(c *gin.Context) {
		api.GetClashSubscription(c)
	})

	subEngine.GET(config.SubPath+"/:token/surge", func(c *gin.Context) {
		api.GetSurgeSubscription(c)
	})

	return runEngine(subEngine, config.Domain, config.SubPort)
}

func runEngine(r *gin.Engine, domain string, port int) error {
	addr := fmt.Sprintf(":%d", port)
	if domain != "" {
		var cert database.Cert
		database.DB.Where("domain = ?", domain).First(&cert)
		if cert.CertPath == "" || cert.KeyPath == "" {
			return fmt.Errorf("域名 %s 证书路径未配置", domain)
		}
		go r.RunTLS(addr, cert.CertPath, cert.KeyPath)
	} else {
		go r.Run(addr)
	}
	return nil
}
