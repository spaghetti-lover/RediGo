package config

import (
	"log"
	"os"
	"strconv"
)

// Redis server configuration
var (
	Protocol      = getEnv("REDIS_PROTOCOL", "tcp")
	Port          = getEnv("REDIS_PORT", ":6379")
	MaxConnection = getEnvAsInt("REDIS_MAX_CONNECTION", 20000)
)

// HTTP Gateway configuration
var (
	HTTPPort         = getEnv("HTTP_PORT", ":8080")
	HTTPReadTimeout  = getEnvAsInt("HTTP_READ_TIMEOUT", 15)
	HTTPWriteTimeout = getEnvAsInt("HTTP_WRITE_TIMEOUT", 15)
)

// Gateway to Redis connection configuration
var (
	RedisConnTimeout = getEnvAsInt("REDIS_CONN_TIMEOUT", 5)
	RedisRWTimeout   = getEnvAsInt("REDIS_RW_TIMEOUT", 10)
	RedisAddr        = getEnv("REDIS_ADDR", "localhost:6379")
)

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}
