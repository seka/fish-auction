package config

import (
	"log/slog"
	"os"
	"strconv"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetLogLevel は LOG_LEVEL を返す。logger.Init を config 構造体の Load より前に走らせるため、
// LogLevel だけは独立した取得経路を用意している。
func GetLogLevel() string {
	return GetEnv("LOG_LEVEL", "info")
}

// GetEnvInt retrieves the value of the environment variable named by the key as an integer.
func GetEnvInt(key string, defaultValue int) int {
	valueStr := GetEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		slog.Warn("invalid env value, using default", "key", key, "value", valueStr, "default", defaultValue)
		return defaultValue
	}
	return value
}
