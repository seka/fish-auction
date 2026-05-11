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
	}
}

func (c *InitAdminConfig) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSslMode)
}
