package config

import (
	"fmt"
	"os"
)

// TestConfig はテスト用の設定
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
}

// LoadTest はテスト用の設定を読み込む
func LoadTest() *TestConfig {
	return &TestConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
	}
}

// AdminConnStr は管理用 DB の接続文字列を返す
func (c *TestConfig) AdminConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword)
}

// TestDBConnStr はテスト用 DB の接続文字列を返す
func (c *TestConfig) TestDBConnStr(dbName string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, dbName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
