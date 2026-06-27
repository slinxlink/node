package service

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

func hopChainName(port int) string {
	return fmt.Sprintf("hop_%d", port)
}

func nft(rule string) error {
	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(rule)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nftables 错误: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func getExcludedPorts() []int {
	var config database.Config
	database.DB.First(&config)

	excluded := []int{config.Port, config.SubPort}

	for p := range util.ReservedPorts {
		excluded = append(excluded, p)
	}

	var inbounds []database.Inbound
	database.DB.Find(&inbounds)
	for _, ib := range inbounds {
		if !ib.HopEnabled {
			excluded = append(excluded, ib.Port)
		}
	}

	return excluded
}

func EnableHop() error {
	var inbounds []database.Inbound
	database.DB.Where("enable = ? AND hop_enabled = ?", true, true).Find(&inbounds)

	if len(inbounds) == 0 {
		return nil
	}

	excluded := getExcludedPorts()
	ports := make([]string, len(excluded))
	for i, p := range excluded {
		ports[i] = fmt.Sprintf("%d", p)
	}

	if err := nft("add table ip slinx_hop"); err != nil {
		return err
	}

	for _, ib := range inbounds {
		parts := strings.SplitN(ib.HopPort, "-", 2)
		if len(parts) != 2 {
			util.Warn("[hop] 端口范围格式错误: %s，跳过入站 %d", ib.HopPort, ib.Port)
			continue
		}

		chain := hopChainName(ib.Port)

		if err := nft(fmt.Sprintf("add chain ip slinx_hop %s { type nat hook prerouting priority -100; }", chain)); err != nil {
			return err
		}
		if err := nft(fmt.Sprintf("add rule ip slinx_hop %s udp dport %s-%s udp dport != { %s } dnat to :%d",
			chain, parts[0], parts[1], strings.Join(ports, ", "), ib.Port)); err != nil {
			return err
		}

		util.Info("[hop] 端口跳跃已启动: %s → %d", ib.HopPort, ib.Port)
	}

	return nil
}

func DisableHop() error {
	return nft("flush table ip slinx_hop")
}

func ValidateHopPort(hopPort string, excludeInboundID uint) string {
	if hopPort == "" {
		return "请填写端口跳跃范围"
	}
	parts := strings.SplitN(hopPort, "-", 2)
	if len(parts) != 2 {
		return "端口跳跃范围格式错误，应为 1000-2000"
	}
	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return "端口跳跃范围必须为数字"
	}
	if start < 1 || start > 65535 || end < 1 || end > 65535 {
		return "端口跳跃范围必须在 1-65535 之间"
	}
	if end <= start {
		return "结束端口必须大于起始端口"
	}

	var others []database.Inbound
	database.DB.Where("id != ? AND hop_enabled = ? AND enable = ?", excludeInboundID, true, true).Find(&others)
	for _, o := range others {
		oParts := strings.SplitN(o.HopPort, "-", 2)
		if len(oParts) != 2 {
			continue
		}
		oStart, _ := strconv.Atoi(strings.TrimSpace(oParts[0]))
		oEnd, _ := strconv.Atoi(strings.TrimSpace(oParts[1]))
		if start <= oEnd && end >= oStart {
			return fmt.Sprintf("端口范围与入站 %d (%s) 重叠", o.Port, o.HopPort)
		}
	}
	return ""
}

func ValidateHopInterval(hopInterval string) string {
	if hopInterval == "" {
		return "请填写端口跳跃间隔"
	}
	parts := strings.SplitN(hopInterval, "-", 2)
	if len(parts) != 2 {
		return "跳跃间隔格式错误，应为 5-10"
	}
	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return "跳跃间隔必须为数字"
	}
	if start < 5 {
		return "跳跃间隔最小为 5 秒"
	}
	if end <= start {
		return "结束间隔必须大于起始间隔"
	}
	return ""
}
