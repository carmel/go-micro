package logger

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrMissingValue = errors.New("(MISSING)")

type zapSugaredLogger struct {
	*zap.SugaredLogger
	// fields []interface{}
}

func NewZapSugaredLogger(opt Options) (*zapSugaredLogger, error) {

	// if len(fields) == 0 {
	// 	fields = []any{"source", Caller(2)}
	// } else {
	// 	fields = append(fields, "source", Caller(2))
	// }

	var writer zapcore.WriteSyncer
	var enc zapcore.Encoder
	if opt.LogLevel == DEBUG {

		writer = zapcore.AddSync(os.Stdout)
		encoder := zap.NewDevelopmentEncoderConfig()
		encoder.EncodeTime = zapcore.TimeEncoderOfLayout("15:05:05.000")
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// encoder.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		// 	enc.AppendString(l.CapitalString())
		// 	switch l {
		// 	case zapcore.DebugLevel:
		// 		enc.AppendString(color.MagentaString("DEBUG"))
		// 	case zapcore.InfoLevel:
		// 		enc.AppendString(color.BlueString("INFO"))
		// 	case zapcore.WarnLevel:
		// 		enc.AppendString(color.YellowString("INFO"))
		// 	case zapcore.ErrorLevel:
		// 		enc.AppendString(color.RedString("INFO"))
		// 	}
		// }

		enc = zapcore.NewConsoleEncoder(encoder)

	} else {
		logWriter, err := NewLogWriter(opt)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(logWriter)

		encoder := zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05")

		// encoder.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		// 	enc.AppendString(strings.Join([]string{caller.TrimmedPath(), runtime.FuncForPC(caller.PC).Name()}, ":"))
		// }

		encoder.TimeKey = "ts"
		encoder.LevelKey = "level"
		encoder.NameKey = "logger"
		encoder.CallerKey = "line"
		encoder.MessageKey = "msg"
		encoder.StacktraceKey = "stacktrace"

		enc = zapcore.NewJSONEncoder(encoder)
	}

	////////////////////////
	// consoleWriter := zapcore.Lock(os.Stdout)

	// core := zapcore.NewTee(
	// 	zapcore.NewCore(enc, zapcore.NewJSONEncoder(cfg.EncoderConfig), cfg.Level),
	// 	zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), consoleWriter, cfg.Level),
	// )
	////////////////////////

	zcore := zapcore.NewCore(
		enc,                                 // 编码器配置
		zapcore.NewMultiWriteSyncer(writer), // 打印日志
		zap.NewAtomicLevelAt(zapcore.Level(opt.LogLevel)), // 日志级别
	)

	hostname, _ := os.Hostname()

	return &zapSugaredLogger{
		SugaredLogger: zap.New(
			zcore,
			// zap.AddStacktrace(zapcore.ErrorLevel),
			// zap.WithCaller(true),
			zap.AddCaller(),
			zap.AddCallerSkip(2),
		).Sugar().With(
			"pid", os.Getegid(),
			"hostname", hostname,
			"go_version", runtime.Version(),
		),
	}, nil

}

func (zl zapSugaredLogger) With(kv ...interface{}) Logger {

	// kvs := append(zl.fields, kv...)
	// n := len(kv)
	// if n != 0 && n%2 != 0 {
	// 	return zl, ErrMissingValue
	// }

	// fld := make([]interface{}, 0, n/2)
	// for i := 0; i < n; i += 2 {
	// 	fld = append(fld, zap.Any(fmt.Sprint(kv[i]), kv[i+1]))
	// }

	zl.SugaredLogger = zl.SugaredLogger.With(kv...)
	return zl

}

func (zl zapSugaredLogger) Log(l Level, msg string) {
	switch l {
	case INFO:
		zl.Info(msg)
	case DEBUG:
		zl.Debug(msg)
	case WARN:
		zl.Warn(msg)
	case ERROR:
		zl.Error(msg)
	}
}

func (zl zapSugaredLogger) Logf(l Level, format string, args ...interface{}) {
	zl.Log(l, fmt.Sprintf(format, args...))
}
