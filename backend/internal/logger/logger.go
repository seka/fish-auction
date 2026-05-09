// Package logger initializes the application-wide structured logger.
// CloudWatch Logs Insights で扱いやすい JSON 構造化ログを既定にする。
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init configures slog as the default logger with a JSON handler writing to stderr.
// LOG_LEVEL 環境変数（debug/info/warn/error）が指定されていれば優先し、未指定時は
// 引数で渡された level を使う。標準ライブラリ log のフォールバック先も差し替えるため、
// log.Print* 経由の呼び出しも slog 経由で出力される。
func Init(level slog.Level) {
	if envLevel, ok := levelFromEnv(); ok {
		level = envLevel
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func levelFromEnv() (slog.Level, bool) {
	raw := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if raw == "" {
		return slog.LevelInfo, false
	}
	switch strings.ToLower(raw) {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}
