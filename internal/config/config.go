package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the HelloAI-AVS
type Config struct {
	// Together AI Configuration
	TogetherAPIKey string
	DefaultModel   string
	MaxTokens      int
	Temperature    float64

	// Request Configuration
	TimeoutSeconds int
	RetryAttempts  int

	// Logging
	LogLevel string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		TogetherAPIKey: getEnv("TOGETHER_API_KEY", ""),
		DefaultModel:   getEnv("DEFAULT_MODEL", "meta-llama/Llama-2-7b-chat-hf"),
		MaxTokens:      getEnvAsInt("MAX_TOKENS", 150),
		Temperature:    getEnvAsFloat("TEMPERATURE", 0.7),
		TimeoutSeconds: getEnvAsInt("TIMEOUT_SECONDS", 30),
		RetryAttempts:  getEnvAsInt("RETRY_ATTEMPTS", 3),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	return config, nil
}

// GetTimeout returns timeout as time.Duration
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

// Helper functions

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsFloat(name string, defaultVal float64) float64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
} 