package logger

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true, ReplaceAttr: replacer})))
}

var replacer = func(groups []string, a slog.Attr) slog.Attr {

	switch a.Key {
	case slog.TimeKey:
		a.Value = slog.StringValue(a.Value.Time().Format("15:05:05.000"))
	// case slog.LevelKey:
	// 	switch a.Value.String() {
	// 	case slog.LevelDebug.String():
	// 		a.Value = slog.StringValue(color.MagentaString("DEBUG"))
	// 	case slog.LevelInfo.String():
	// 		a.Value = slog.StringValue(color.BlueString("INFO"))
	// 	case slog.LevelWarn.String():
	// 		a.Value = slog.StringValue(color.YellowString("WARN"))
	// 	case slog.LevelError.String():
	// 		a.Value = slog.AnyValue(color.RedString("ERROR"))
	// 	}
	// case slog.MessageKey:
	// 	a.Value = slog.StringValue(color.CyanString(a.Value.String()))
	case slog.SourceKey:
		_, file, line, _ := runtime.Caller(7)
		a.Value = slog.StringValue(Caller(file, line))

	}
	return a

}

func Debugf(format string, arg ...any) {
	slog.Debug(fmt.Sprintf(format, arg...))
}

func Infof(format string, arg ...any) {
	slog.Info(fmt.Sprintf(format, arg...))
}

func Warnf(format string, arg ...any) {
	slog.Warn(fmt.Sprintf(format, arg...))
}

func Errorf(format string, arg ...any) {
	slog.Error(fmt.Sprintf(format, arg...))
}

func Log(l Level, msg string) {
	switch l {
	case INFO:
		slog.Info(msg)
	case DEBUG:
		slog.Debug(msg)
	case WARN:
		slog.Warn(msg)
	case ERROR:
		slog.Error(msg)
	}
}

func With(args ...any) *slog.Logger {
	return slog.With(args...)
}
