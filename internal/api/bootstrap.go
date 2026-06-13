package api

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/seekky/slinx-node/internal/core"
	"github.com/seekky/slinx-node/internal/util"
)

func Restart(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
	go restart()
}

func restart() {
	util.Info("[bootstrap] 面板重启")
	core.Default.Stop()
	exe, _ := os.Executable()
	syscall.Exec(exe, os.Args, os.Environ())
}
