package config

import "fmt"

// InitAdminConfig represents the configuration for the init_admin CLI.
type InitAdminConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSslMode  string
	AppEnv           string
}

// NewInitAdminConfig loads configuration for the init_admin command.
func NewInitAdminConfig() *InitAdminConfig {
	return &InitAdminConfig{
		PostgresHost:     GetEnv("POSTGRES_HOST", ""),
		PostgresPort:     GetEnv("POSTGRES_PORT", ""),
		PostgresUser:     GetEnv("POSTGRES_USER", ""),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", ""),
		PostgresSslMode:  GetEnv("POSTGRES_SSLMODE", "disable"),
		AppEnv:           GetEnv("APP_ENV", "development"),
	}
}

func (c *InitAdminConfig) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSslMode)
}

// Validate は DB 接続に必須な値の不足を起動前に検出する。
// NewInitAdminConfig は空文字をデフォルトに用いるため、未設定時に DB 接続段階で
// 失敗するまで原因が分かりづらいことを避ける。
func (c *InitAdminConfig) Validate() error {
	var missing []string
	if c.PostgresHost == "" {
		missing = append(missing, "POSTGRES_HOST")
	}
	if c.PostgresPort == "" {
		missing = append(missing, "POSTGRES_PORT")
	}
	if c.PostgresUser == "" {
		missing = append(missing, "POSTGRES_USER")
	}
	if c.PostgresPassword == "" {
		missing = append(missing, "POSTGRES_PASSWORD")
	}
	if c.PostgresDB == "" {
		missing = append(missing, "POSTGRES_DB")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", missing)
	}
	if err := validateSSLMode(c.AppEnv, c.PostgresSslMode); err != nil {
		return err
	}
	return nil
}
