package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want slog.Level
	}{
		{name: "empty falls back to info", raw: "", want: slog.LevelInfo},
		{name: "unknown falls back to info", raw: "verbose", want: slog.LevelInfo},
		{name: "lowercase debug", raw: "debug", want: slog.LevelDebug},
		{name: "uppercase DEBUG", raw: "DEBUG", want: slog.LevelDebug},
		{name: "lowercase info", raw: "info", want: slog.LevelInfo},
		{name: "uppercase INFO", raw: "INFO", want: slog.LevelInfo},
		{name: "lowercase warn", raw: "warn", want: slog.LevelWarn},
		{name: "warning alias", raw: "warning", want: slog.LevelWarn},
		{name: "uppercase WARNING", raw: "WARNING", want: slog.LevelWarn},
		{name: "lowercase error", raw: "error", want: slog.LevelError},
		{name: "uppercase ERROR", raw: "ERROR", want: slog.LevelError},
		{name: "trims whitespace", raw: "  debug  ", want: slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseLevel(tt.raw))
		})
	}
}
