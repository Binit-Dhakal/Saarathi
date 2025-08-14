package log

import (
	"context"
	"log/slog"
	"os"
)

type ILogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, err error, args ...any)
}

type Logger struct {
	logger *slog.Logger
}

func NewStandardLogger() ILogger {
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		logger: slog.New(stdoutHandler),
	}
}

func (l *Logger) log(level slog.Level, msg string, args ...any) {
	l.logger.Log(context.TODO(), level, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, err error, args ...any) {
	args = append(args, "error", err)
	l.log(slog.LevelError, msg, args...)
}
