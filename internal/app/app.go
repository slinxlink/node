package app

import (
	"os"
	"path/filepath"
)

var Version = "dev"

var Dir = func() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}()
