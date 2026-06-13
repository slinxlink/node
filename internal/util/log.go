package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *log.Logger
var currentLevel = "info"
var enabled = true
var levelMap = map[string]int{"info": 0, "warn": 1, "error": 2}

func InitLog(path, level string, enable bool) {
	writer := io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename: path,
			MaxAge:   1,
			MaxSize:  10,
			Compress: false,
		},
	)
	logger = log.New(writer, "", 0)
	currentLevel = level
	enabled = enable
}

func Info(format string, args ...any) {
	write("INFO", format, args...)
}

func Warn(format string, args ...any) {
	write("WARN", format, args...)
}

func Error(format string, args ...any) {
	write("ERROR", format, args...)
}

func write(level, format string, args ...any) {
	if !enabled {
		return
	}
	if levelMap[strings.ToLower(level)] < levelMap[strings.ToLower(currentLevel)] {
		return
	}
	msg := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), level, fmt.Sprintf(format, args...))
	if logger != nil {
		logger.Println(msg)
	} else {
		fmt.Println(msg)
	}
}
