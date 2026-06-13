package util

import (
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
)

func GenerateString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GeneratePort() int {
	return rand.Intn(55535) + 10000
}

var realityTargets = []string{
	"www.microsoft.com",
	"www.amazon.com",
	"aws.amazon.com",
	"www.cloudflare.com",
	"www.intel.com",
	"www.sony.com",
	"www.samsung.com",
	"www.nvidia.com",
	"www.amd.com",
	"www.github.com",
	"addons.mozilla.org",
}

func GenerateRealityTarget() string {
	return realityTargets[rand.Intn(len(realityTargets))]
}

// 生成单个 ShortID，指定长度（字节数，输出为 2*n 位 hex）
func GenerateShortID(bytes int) (string, error) {
	out, err := exec.Command("bin/sing-box", "generate", "rand", "--hex", strconv.Itoa(bytes)).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// 生成 Reality 密钥对
func GenerateRealityKeyPair() (privateKey, publicKey string, err error) {
	out, err := exec.Command("bin/sing-box", "generate", "reality-keypair").Output()
	if err != nil {
		return "", "", err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PrivateKey:") {
			privateKey = strings.TrimSpace(strings.TrimPrefix(line, "PrivateKey:"))
		}
		if strings.HasPrefix(line, "PublicKey:") {
			publicKey = strings.TrimSpace(strings.TrimPrefix(line, "PublicKey:"))
		}
	}
	return
}

// 生成8个 ShortID，第一个长度随机但不短于4字节（8位hex），其余完全随机
func GenerateShortIDs() ([]string, error) {
	var ids []string
	for i := 0; i < 8; i++ {
		var length int
		if i == 0 {
			length = rand.Intn(6) + 3 // 3~8字节，6~16位hex
		} else {
			length = rand.Intn(8) + 1 // 1~8字节，2~16位hex
		}
		id, err := GenerateShortID(length)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// 生成 UUID
func GenerateUUID() (string, error) {
	out, err := exec.Command("bin/sing-box", "generate", "uuid").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
