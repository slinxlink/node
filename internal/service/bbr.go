package service

import (
	"os"
	"os/exec"
	"strings"

	"github.com/slinxlink/node/internal/database"
)

const bbrConf = "/etc/sysctl.d/99-bbr.conf"

func BBREnable() error {
	content := "net.core.default_qdisc=fq\nnet.ipv4.tcp_congestion_control=bbr\n"
	if err := os.WriteFile(bbrConf, []byte(content), 0644); err != nil {
		return err
	}
	if err := exec.Command("sysctl", "--system").Run(); err != nil {
		return err
	}
	exec.Command("sysctl", "-w", "net.core.default_qdisc=fq").Run()
	exec.Command("sysctl", "-w", "net.ipv4.tcp_congestion_control=bbr").Run()
	return nil
}

func BBRDisable() error {
	os.Remove(bbrConf)
	content := "net.core.default_qdisc=pfifo_fast\nnet.ipv4.tcp_congestion_control=cubic\n"
	if err := os.WriteFile("/etc/sysctl.d/99-bbr-disable.conf", []byte(content), 0644); err != nil {
		return err
	}
	if err := exec.Command("sysctl", "--system").Run(); err != nil {
		return err
	}
	exec.Command("sysctl", "-w", "net.core.default_qdisc=pfifo_fast").Run()
	exec.Command("sysctl", "-w", "net.ipv4.tcp_congestion_control=cubic").Run()
	return nil
}

func BBRStatus() bool {
	out, err := exec.Command("sysctl", "net.ipv4.tcp_congestion_control").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "bbr")
}

func BBRApply(enable bool) bool {
	if enable {
		BBREnable()
	} else {
		BBRDisable()
	}
	actual := BBRStatus()
	database.DB.Model(&database.Config{}).Where("id = 1").Update("bbr", actual)
	return actual
}
