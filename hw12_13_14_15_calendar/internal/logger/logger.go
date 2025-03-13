package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

type Logger struct {
	handler *slog.Logger
}

func (l *Logger) Warn(msg string)  { l.handler.Warn(msg) }
func (l *Logger) Info(msg string)  { l.handler.Info(msg) }
func (l *Logger) Error(msg string) { l.handler.Error(msg) }

func New(level string) *Logger {
	var sl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		sl = slog.LevelDebug
	case "info":
		sl = slog.LevelInfo
	case "warn":
		sl = slog.LevelWarn
	case "error":
		sl = slog.LevelError
	default:
		sl = slog.LevelInfo
	}

	logConfig := &slog.HandlerOptions{
		Level:       sl,
		ReplaceAttr: replaceAttr,
	}
	logHandler := slog.NewJSONHandler(os.Stdout, logConfig)

	return &Logger{
		handler: slog.New(logHandler),
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	_ = groups
	switch a.Key {
	case slog.TimeKey:
		return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
	case slog.MessageKey:
		return slog.String("rest", a.Value.String())
	case slog.LevelKey:
		return slog.String("severity", a.Value.String())
	default:
		return a
	}
}
