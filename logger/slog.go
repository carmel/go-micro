package logger

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

type Slogger struct {
	*slog.Logger
}

func (l Level) Level() slog.Level {
	return slog.Level(l * 4)
}

func NewSlogger(opt Options) Slogger {
	hostname, _ := os.Hostname()
	if opt.LogLevel == DEBUG {

		return Slogger{
			slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						AddSource:   true,
						Level:       opt.LogLevel,
						ReplaceAttr: replacer,
					},
				).WithAttrs([]slog.Attr{
					slog.Group("app",
						slog.Int("pid", os.Getpid()),
						slog.String("hostname", hostname),
						slog.String("go_version", runtime.Version()),
					),
				}),
			),
		}
	} else {
		return Slogger{
			slog.New(
				slog.NewJSONHandler(
					os.Stdout,
					&slog.HandlerOptions{
						AddSource:   true,
						Level:       opt.LogLevel,
						ReplaceAttr: replacer,
					},
				).WithAttrs([]slog.Attr{
					slog.Group("app",
						slog.Int("pid", os.Getpid()),
						slog.String("hostname", hostname),
						slog.String("go_version", runtime.Version()),
					),
				}),
			),
		}
	}
}

func (s Slogger) With(kv ...interface{}) Logger {
	s.Logger = s.Logger.With(kv...)
	return s
}

func (s Slogger) Log(l Level, msg string) {
	switch l {
	case INFO:
		s.Info(msg)
	case DEBUG:
		s.Debug(msg)
	case WARN:
		s.Warn(msg)
	case ERROR:
		s.Error(msg)
	}
}

func (s Slogger) Logf(l Level, format string, args ...interface{}) {
	s.Log(l, fmt.Sprintf(format, args...))
}

// type ctxTraceIdKey struct{}

// type ContextHandler struct {
// 	handler slog.Handler
// }

// func NewContextHandler(h slog.Handler) *ContextHandler {
// 	if lh, ok := h.(*ContextHandler); ok {
// 		h = lh.handler
// 	}
// 	return &ContextHandler{h}
// }

// func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
// 	return h.handler.Enabled(ctx, level)
// }

// func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
// 	if ctx == nil {
// 		return h.handler.Handle(ctx, r)
// 	}
// 	if tid, ok := ctx.Value(ctxTraceIdKey{}).(string); ok {
// 		traceAttr := slog.Attr{
// 			Key:   "trace_id",
// 			Value: slog.StringValue(tid),
// 		}
// 		r.AddAttrs(traceAttr)
// 	}
// 	return h.handler.Handle(ctx, r)
// }

// func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
// 	return NewContextHandler(h.handler.WithAttrs(attrs))
// }

// func (h *ContextHandler) WithGroup(name string) slog.Handler {
// 	return NewContextHandler(h.handler.WithGroup(name))
// }
