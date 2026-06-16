package util

import (
	"regexp"
	"strings"
)

var unsafeFileName = regexp.MustCompile(`[\x00-\x1f\\/:*?"<>|]`)

// SanitizeFileName 清除文件名中的非法字符
func SanitizeFileName(name string) string {
	name = unsafeFileName.ReplaceAllString(name, "_")
	return strings.Trim(name, ". ")
}
