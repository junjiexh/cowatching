package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	ServerPort        string
	AWSRegion         string
	AWSAccessKeyID    string
	AWSSecretAccessKey string
	S3BucketName      string
	S3VideoPrefix     string
}

func Load() (*Config, error) {
	// Load .env file if it exists (optional in production)
	_ = godotenv.Load()

	config := &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "cowatching"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID:    getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		S3BucketName:      getEnv("S3_BUCKET_NAME", ""),
		S3VideoPrefix:     getEnv("S3_VIDEO_PREFIX", "videos/"),
	}

	return config, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
