// Package logger initializes the application-wide structured logger.
// CloudWatch Logs Insights で扱いやすい JSON 構造化ログを既定にする。
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init configures slog as the default logger with a JSON handler writing to stderr.
// level は config 経由で渡される文字列（debug/info/warn/error）。空文字や未知値は info にフォールバックする。
func Init(level string) {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     parseLevel(level),
	})
	slog.SetDefault(slog.New(handler))
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
