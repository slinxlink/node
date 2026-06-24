package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

type Manager struct {
	cmd        *exec.Cmd
	BinPath    string
	ConfigPath string
	running    bool
	mu         sync.Mutex
	waitDone   chan struct{}
}

var Default = &Manager{
	BinPath:    "bin/sing-box",
	ConfigPath: "data/sing-box.json",
}

func (m *Manager) Init() {
	var coreRecord database.Core
	database.DB.First(&coreRecord)
	m.BinPath = coreRecord.BinPath
	m.ConfigPath = coreRecord.ConfigPath
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}
	if _, err := os.Stat(m.ConfigPath); os.IsNotExist(err) {
		if err := os.WriteFile(m.ConfigPath, []byte("{}"), 0644); err != nil {
			return err
		}
	}

	m.cmd = exec.Command(m.BinPath, "run", "-c", m.ConfigPath)
	stderr, _ := m.cmd.StderrPipe()

	if err := m.cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			util.Warn("[core] %s", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			util.Warn("[core] stderr读取错误: %v", err)
		}
	}()

	m.waitDone = make(chan struct{})
	m.running = true
	util.Info("[core] 核心启动成功")
	go m.watch()
	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()

	if m.cmd == nil || !m.running {
		m.mu.Unlock()
		return nil
	}

	cmd := m.cmd
	waitDone := m.waitDone
	m.running = false
	m.cmd = nil
	m.mu.Unlock()

	cmd.Process.Signal(syscall.SIGTERM)

	select {
	case <-waitDone:
		// 优雅退出完成
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		<-waitDone
	}

	exec.Command("pkill", "-f", "sing-box").Run()
	util.Info("[core] 核心已停止")
	return nil
}

func (m *Manager) Restart() error {
	if err := m.Stop(); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	util.Info("[core] 核心重启")
	return m.Start()
}

func (m *Manager) Check() error {
	out, err := exec.Command(m.BinPath, "check", "-c", m.ConfigPath).CombinedOutput()
	if err != nil {
		return errors.New(strings.TrimSpace(string(out)))
	}
	return nil
}

func (m *Manager) Apply() error {
	if err := generateConfig(); err != nil {
		return err
	}
	if err := m.Check(); err != nil {
		util.Error("[core] 配置校验失败: %v", err)
		return err
	}
	if m.Status() == "running" {
		return m.Restart()
	}
	return nil
}

func (m *Manager) Status() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return "running"
	}
	return "stopped"
}

func (m *Manager) Version() string {
	out, err := exec.Command(m.BinPath, "version").Output()
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return parts[2]
		}
	}
	return "unknown"
}

func (m *Manager) watch() {
	cmd := m.cmd
	waitDone := m.waitDone
	cmd.Wait()
	close(waitDone)

	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return
	}
	m.running = false
	m.cmd = nil
	util.Warn("[core] 核心意外退出")
}

func (m *Manager) Process() (map[string]any, error) {
	if !m.running || m.cmd == nil {
		return nil, fmt.Errorf("核心未运行")
	}

	resp, err := http.Get("http://127.0.0.1:9090/memory")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mem struct {
		Inuse uint64 `json:"inuse"`
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &mem); err == nil && mem.Inuse > 0 {
			break
		}
	}
	if scanner.Err() != nil {
		util.Warn("[core] 读取内存数据失败: %v", scanner.Err())
	}

	pid := m.cmd.Process.Pid

	return map[string]any{
		"memory":  mem.Inuse,
		"threads": getThread(pid),
	}, nil
}
