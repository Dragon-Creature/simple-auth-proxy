package env

import "os"

func GetEnvOrDefault(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}
	return value
}
