package core

import (
	"os"
)

func (m *Manager) ReadConfig() ([]byte, error) {
	return os.ReadFile(m.ConfigPath)
}

func (m *Manager) WriteConfig(data []byte) error {
	return os.WriteFile(m.ConfigPath, data, 0644)
}
