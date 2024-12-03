package utils

import "os"

// GetEnv is a helper function to get environment variables with a fallback value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
