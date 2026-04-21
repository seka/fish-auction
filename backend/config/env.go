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

func loadFrontendURL() (*url.URL, error) {
	frontendURLStr := GetEnv("FRONTEND_URL", "https://localhost")
	frontendURL, err := url.Parse(frontendURLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid FRONTEND_URL: %w", err)
	}
	if frontendURL.Scheme == "" || frontendURL.Host == "" {
		return nil, fmt.Errorf("invalid FRONTEND_URL: missing scheme or host")
	}
	return frontendURL, nil
}
