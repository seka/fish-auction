package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves the value of the environment variable named by the key as an integer.
func GetEnvInt(key string, defaultValue int) int {
	valueStr := GetEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

func loadFrontendURL() *url.URL {
	frontendURLStr := GetEnv("FRONTEND_URL", "https://localhost")
	frontendURL, err := url.Parse(frontendURLStr)
	if err != nil || frontendURL.Scheme == "" || frontendURL.Host == "" {
		return nil
	}
	return frontendURL
}

func validateBase(host, port, user, password, db string, frontendURL *url.URL) error {
	if host == "" || port == "" || user == "" || password == "" || db == "" {
		return fmt.Errorf("missing required environment variables (POSTGRES_*)")
	}
	if frontendURL == nil {
		return fmt.Errorf("invalid or missing FRONTEND_URL")
	}
	return nil
}
