package main

import (
	"embed"
	"os"

	"github.com/seekky/slinx-node/internal/bootstrap"
	"github.com/seekky/slinx-node/internal/cli"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	os.Chdir("var")

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		cli.Start()
		return
	}
	bootstrap.Start(webFS)
}
