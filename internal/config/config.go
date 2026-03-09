package config

import (
	"os"
	"strings"
)

type Config struct {
	Port          string
	AllowedOrgins []string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		AllowedOrgins: getEnvSlice("ALLOWED_ORIGINS",
			[]string{"http://localhost:3000"}),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "rohit"),
		DBPassword: getEnv("DB_PASSWORD", "learngo123"),
		DBName:     getEnv("DB_NAME", "goframework"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	if val := os.Getenv(key); val != "" {
		origins := strings.Split(val, ",")
		return origins
	}

	return fallback
}
