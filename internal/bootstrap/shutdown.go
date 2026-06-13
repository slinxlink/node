package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/slinxlink/node/internal/core"
)

func Shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	core.Default.Stop()
	os.Exit(0)
}
