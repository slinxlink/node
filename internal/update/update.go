package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

type CheckResult struct {
	HasUpdate     bool   `json:"has_update"`
	LatestVersion string `json:"latest_version"`
	Changelog     string `json:"changelog"`
}

type releaseInfo struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

func getRepo() string {
	var config database.Config
	database.DB.First(&config)
	return config.Repo
}

func apiURL() string {
	path := strings.TrimPrefix(getRepo(), "https://github.com/")
	return "https://api.github.com/repos/" + path + "/releases/latest"
}

func downloadURL(version string) string {
	return fmt.Sprintf("%s/releases/download/%s/slinx_linux_%s", getRepo(), version, runtime.GOARCH)
}

func Check(currentVersion string) (*CheckResult, error) {
	resp, err := http.Get(apiURL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release releaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &CheckResult{
		HasUpdate:     release.TagName != "" && release.TagName != currentVersion,
		LatestVersion: release.TagName,
		Changelog:     release.Body,
	}, nil
}

func Update(currentVersion string) error {
	result, err := Check(currentVersion)
	if err != nil {
		util.Error("[update] 版本检测失败: %v", err)
		return err
	}
	if !result.HasUpdate {
		return nil
	}

	util.Info("[update] 开始更新: %s → %s", currentVersion, result.LatestVersion)

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	tmpDir := filepath.Join(filepath.Dir(exe), "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tmp := filepath.Join(tmpDir, "slinx.tmp")
	if err := exec.Command("wget", "-q", "-O", tmp, downloadURL(result.LatestVersion)).Run(); err != nil {
		util.Error("[update] 下载失败: %v", err)
		return fmt.Errorf("下载失败: %w", err)
	}

	if err := os.Rename(tmp, exe); err != nil {
		util.Error("[update] 替换文件失败: %v", err)
		return err
	}

	if err := os.Chmod(exe, 0755); err != nil {
		util.Error("[update] 设置权限失败: %v", err)
		return err
	}

	util.Info("[update] 更新成功: %s，面板即将重启", result.LatestVersion)
	exec.Command("systemctl", "restart", "slinx").Run()
	return nil
}
