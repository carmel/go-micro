package log

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/microservices/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrMissingValue = errors.New("(MISSING)")

type zapSugaredLogger struct {
	preload []interface{}
	*zap.SugaredLogger
}

func NewZapSugaredLogger(opt *Options) (*zapSugaredLogger, error) {

	var encoder = zap.NewProductionEncoderConfig()

	encoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(util.FORMAT_ISO8601_DATE_TIME_MILLI))
	}

	// encoder.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// 	enc.AppendString(strings.Join([]string{caller.TrimmedPath(), runtime.FuncForPC(caller.PC).Name()}, ":"))
	// }

	encoder.TimeKey = "ts"
	encoder.LevelKey = "level"
	encoder.NameKey = "logger"
	encoder.CallerKey = "line"
	encoder.MessageKey = "msg"
	encoder.StacktraceKey = "stacktrace"

	////////////////////////
	// consoleWriter := zapcore.Lock(os.Stdout)

	// core := zapcore.NewTee(
	// 	zapcore.NewCore(enc, zapcore.NewJSONEncoder(cfg.EncoderConfig), cfg.Level),
	// 	zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), consoleWriter, cfg.Level),
	// )
	////////////////////////
	logWriter, err := NewLogWriter(opt)

	if err != nil {
		return nil, err
	}

	var logLevel zapcore.Level

	switch opt.LogLevel {
	case INFO:
		logLevel = zapcore.InfoLevel
	case DEBUG:
		logLevel = zapcore.DebugLevel
	case WARN:
		logLevel = zapcore.WarnLevel
	case ERROR:
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.ErrorLevel
	}

	zcore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logWriter)), // 打印到控制台和文件
		zap.NewAtomicLevelAt(logLevel), // 日志级别
	)

	return &zapSugaredLogger{
		SugaredLogger: zap.New(zcore, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar(),
	}, nil

}

func (zl *zapSugaredLogger) Load(kv ...interface{}) error {

	if len(kv) == 0 {
		return nil
	}
	kvs := append(zl.preload, kv...)
	if len(kvs)%2 != 0 {
		return ErrMissingValue
	}

	keylen := len(kvs)

	fld := make([]interface{}, 0, keylen/2)
	for i := 0; i < keylen; i += 2 {
		fld = append(fld, zap.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
	}

	zl.SugaredLogger = zl.SugaredLogger.With(fld...)
	return nil

}

func (zl zapSugaredLogger) Log(l level, args ...interface{}) {
	zl.Load("source", Caller(2))
	switch l {
	case INFO:
		zl.Info(args...)
	case DEBUG:
		zl.Debug(args...)
	case WARN:
		zl.Warn(args...)
	case ERROR:
		zl.Error(args...)
	}
}

func (zl zapSugaredLogger) Logf(l level, format string, args ...interface{}) {
	zl.Load("source", Caller(2))
	switch l {
	case INFO:
		zl.Infof(format, args...)
	case DEBUG:
		zl.Debugf(format, args...)
	case WARN:
		zl.Warnf(format, args...)
	case ERROR:
		zl.Errorf(format, args...)
	}
}
