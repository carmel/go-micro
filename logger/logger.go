package logger

import (
	"strconv"
	"strings"
	"time"
)

type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
)

// Options represents optional behavior you can specify for a new LogWriter.
type Options struct {
	// 日志存放路径
	LogPath string `yaml:"log-path"`
	// 日志保存期限
	MaxAge time.Duration `yaml:"max-age"`
	// 单个日志文件大小，超过则rotate
	MaxSize int64 `yaml:"max-size"`
	// 日志是否压缩
	// Compress bool `yaml:"compress"`
	// 日志级别
	LogLevel Level `yaml:"log-level"`
}

type Logger interface {
	With(kv ...interface{}) Logger
	Log(l Level, msg string)
	Logf(l Level, format string, args ...interface{})
}

func Caller(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
	idx = strings.LastIndexByte(file[:idx], '/')
	return file[idx+1:] + ":" + strconv.Itoa(line)
}
