package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	GinMode       string
	UploadDir     string
	MaxUploadSize int64
	NATSURL       string
	CORSOrigin    string
	DBMaxConns    int32
	DBMinConns    int32
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
	FromEmail     string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://livehouse:livehouse_dev@localhost:5432/livehouse_dev?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		GinMode:       getEnv("GIN_MODE", "debug"),
		UploadDir:     getEnv("UPLOAD_DIR", "/app/uploads"),
		MaxUploadSize: 10 * 1024 * 1024,
		NATSURL:       getEnv("NATS_URL", "nats://localhost:4222"),
		CORSOrigin:    getEnv("CORS_ORIGIN", "*"),
		DBMaxConns:    int32(getEnvInt("DB_MAX_CONNS", 25)),
		DBMinConns:    int32(getEnvInt("DB_MIN_CONNS", 5)),
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		FromEmail:     getEnv("FROM_EMAIL", "noreply@livehouseaas.com"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}
