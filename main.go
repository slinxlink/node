package main

import (
	"embed"
	"os"

	"github.com/slinxlink/node/internal/app"
	"github.com/slinxlink/node/internal/bootstrap"
	"github.com/slinxlink/node/internal/cli"
)

var Version = "dev"

//go:embed web/dist
var fs embed.FS

func main() {
	os.Chdir("var")

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		cli.Start(Version)
		return
	}

	app.Version = Version
	bootstrap.Boot(fs)
}
