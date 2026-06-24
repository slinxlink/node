package main

import (
	"embed"
	"os"

	"github.com/slinxlink/node/internal/app"
	"github.com/slinxlink/node/internal/bootstrap"
	"github.com/slinxlink/node/internal/cli"
	"github.com/slinxlink/node/internal/server"
)

var Version = "dev"

//go:embed web/dist
var webFS embed.FS

func main() {
	os.Chdir("var")

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		cli.Start(Version)
		return
	}
	app.Version = Version
	server.Init(webFS)
	bootstrap.Start()
}
