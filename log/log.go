package log

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

type level uint8

const (
	INFO level = iota
	DEBUG
	WARN
	ERROR
)

// Options represents optional behavior you can specify for a new LogWriter.
type Options struct {
	// 日志保存期限
	MaxAge time.Duration `yaml:"max-age"`
	// 单个日志文件大小，超过则rotate
	MaxSize int64 `yaml:"max-size"`
	// 日志是否压缩
	Compress bool `yaml:"compress"`
	// 日志级别
	LogLevel level `yaml:"log-level"`
	// 日志存放路径
	LogPath string `yaml:"log-path"`
}

type Logger interface {
	Load(kv ...interface{}) error
	Log(l level, args ...interface{})
	Logf(l level, format string, args ...interface{})
}

func Caller(depth int) interface{} {
	_, file, line, _ := runtime.Caller(depth)
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
	idx = strings.LastIndexByte(file[:idx], '/')
	return file[idx+1:] + ":" + strconv.Itoa(line)
}
