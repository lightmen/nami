package alog

import (
	"context"
	"sync/atomic"
	"time"
)

var defaultLogger atomic.Value

func init() {
	SetDefault(New(newDefaultHandler()))
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func Default() *Logger {
	l := defaultLogger.Load().(*Logger)
	return l
}

type Logger struct {
	handler Handler
}

func New(h Handler) *Logger {
	return &Logger{handler: h}
}

func (l *Logger) Handler() Handler {
	return l.handler
}

func (l *Logger) clone() *Logger {
	c := *l
	return &c
}

func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	return l.Handler().Enabled(ctx, level)
}

func (l *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
	if !l.Enabled(ctx, level) {
		return
	}

	info := CallerInfo(2)
	r := NewRecord(time.Now(), level, msg, info)
	r.BindArgs(args...)

	if ctx == nil {
		ctx = context.Background()
	}

	_ = l.Handler().Handle(ctx, r)
}

func DebugCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelDebug, format, args...)
}

func Debug(format string, args ...any) {
	Default().log(nil, LevelDebug, format, args...)
}

func InfoCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelInfo, format, args...)
}

func Info(format string, args ...any) {
	Default().log(nil, LevelInfo, format, args...)
}

func ErrorCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelError, format, args...)
}

func Error(format string, args ...any) {
	Default().log(nil, LevelError, format, args...)
}

func FatalCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelFatal, format, args...)
}

func Fatal(format string, args ...any) {
	Default().log(nil, LevelFatal, format, args...)
}

// type Logger interface {
// 	Debug(format string, args ...any)
// 	Info(format string, args ...any)
// 	Error(format string, args ...any)
// 	Fatal(format string, args ...any)
// 	Trace(format string, args ...any)
// 	Warn(format string, args ...any)
// }

// type LogLevel uint32

// // 禁用日志输出, 设置日志等级OFF
// const (
// 	TRACE = iota
// 	DEBUG
// 	INFO
// 	NOTICE
// 	WARN
// 	ERROR
// 	FATAL
// 	OFF
// )

// 日志等级

// func (level LogLevel) String() string {
// 	switch level {
// 	case TRACE:
// 		return "TRACE"
// 	case DEBUG:
// 		return "DEBUG"
// 	case INFO:
// 		return "INFO"
// 	case NOTICE:
// 		return "NOTICE"
// 	case WARN:
// 		return "WARN"
// 	case ERROR:
// 		return "ERROR"
// 	case FATAL:
// 		return "FATAL"
// 	default:
// 		fmt.Fprintf(os.Stderr, "[logsdk] level unknown: %d\n", level)
// 		return "UNKNOWN"
// 	}
// }

// var logLevelName = []struct {
// 	name  string
// 	level LogLevel
// }{
// 	{"trace", TRACE},
// 	{"debug", DEBUG},
// 	{"notice", NOTICE},
// 	{"info", INFO},
// 	{"warn", WARN},
// 	{"error", ERROR},
// 	{"fatal", FATAL},
// }

// func logLevelByName(lv string) (level LogLevel) {
// 	level = TRACE
// 	for _, v := range logLevelName {
// 		if lv == v.name {
// 			level = v.level
// 		}
// 	}
// 	return
// }
