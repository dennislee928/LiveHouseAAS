package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	GinMode       string
	UploadDir     string
	MaxUploadSize int64
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://livehouse:livehouse_dev@localhost:5432/livehouse_dev?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		GinMode:       getEnv("GIN_MODE", "debug"),
		UploadDir:     getEnv("UPLOAD_DIR", "/app/uploads"),
		MaxUploadSize: 10 * 1024 * 1024, // 10MB
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
