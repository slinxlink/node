package config

import (
	"os"
	"path/filepath"

	"github.com/slinxlink/node/internal/database"
)

var Config database.Config

var Version = "dev"

var Dir = func() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}()

func Load() error {
	if err := database.DB.First(&Config).Error; err != nil {
		return err
	}
	Config.Dir = Dir
	return nil
}
