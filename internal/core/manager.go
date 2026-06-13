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

	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/util"
)

type Manager struct {
	cmd        *exec.Cmd
	BinPath    string
	ConfigPath string
	running    bool
	mu         sync.Mutex
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

	m.running = true
	util.Info("[core] 核心启动成功")
	go m.watch()
	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd != nil && m.running {
		m.cmd.Process.Signal(syscall.SIGTERM)
		time.Sleep(500 * time.Millisecond)
		m.cmd.Process.Kill()
	}

	// 兜底：杀掉所有 sing-box 进程
	exec.Command("pkill", "-f", "sing-box").Run()

	m.running = false
	m.cmd = nil
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

// Status 返回核心运行状态
func (m *Manager) Status() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return "running"
	}
	return "stopped"
}

// Version 获取核心版本号
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

// watch 监听进程退出，自动更新状态
func (m *Manager) watch() {
	m.cmd.Wait()
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

	// RAM from Clash API
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

	// Threads from /proc
	pid := m.cmd.Process.Pid

	return map[string]any{
		"memory":  mem.Inuse,
		"threads": getThread(pid),
	}, nil
}
