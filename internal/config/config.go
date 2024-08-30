package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     int
	HTTPPort   string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:     getEnv("DB_USER", "admin"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "tabichan"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		HTTPPort:   getEnv("HTTP_PORT", "localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
