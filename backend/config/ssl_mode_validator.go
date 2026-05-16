package config

import "fmt"

func validateSSLMode(appEnv, sslMode string) error {
	if appEnv == "production" && sslMode == "disable" {
		return fmt.Errorf("POSTGRES_SSLMODE=disable is not allowed in production (APP_ENV=%q)", appEnv)
	}
	return nil
}
