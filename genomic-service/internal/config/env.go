package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadEnvVariable() *EnvVariable {
	return &EnvVariable{
		PrivateKey: getEnvOrDefault("PRIVATE_KEY", ""),
	}
}

// populate env variable
func (cfg *Config) SetupEnvVariable() {
	envVariable := loadEnvVariable()
	cfg.WalletSettings.PrivateKey = envVariable.PrivateKey
}
