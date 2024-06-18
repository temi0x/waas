package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RPC_ETHEREUM string
	RPC_SEPOLIA  string
	SqlUsername  string
	SqlPassword  string
	SqlUrl       string
	DbName       string
	CrypteaKey   string
	// BotToken           string
	// UniSatApiKey       string
	// MagicEdenApiKey    string
	// OrdinalNovusApiKey string
	// OKXSecretkey       string
	// OKXAPIKey          string
	// OKXPassPhrase      string
}

// LoadConfig reads configuration from .env file and environment variables.
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file, loading config from environment")
	}

	config := Config{
		// BotToken:           getEnv("BOT_TOKEN", ""),
		// UniSatApiKey:       getEnv("UNISAT_API_KEY", ""),
		// MagicEdenApiKey:    getEnv("MAGIC_EDEN_API_KEY", ""),
		// OrdinalNovusApiKey: getEnv("ORDINAL_NOVUS_API_KEY", ""),
		// OKXSecretkey:       getEnv("OKX_SECRET_KEY", ""),
		// OKXAPIKey:          getEnv("OKX_API_KEY", ""),
		// OKXPassPhrase:      getEnv("OKX_PASSPHRASE", ""),
		RPC_ETHEREUM: getEnv("RPC_ETHEREUM", ""),
		RPC_SEPOLIA:  getEnv("RPC_SEPOLIA", ""),
		SqlUsername:  getEnv("SQL_USERNAME", ""),
		SqlPassword:  getEnv("SQL_PASSWORD", ""),
		SqlUrl:       getEnv("SQL_URL", ""),
		DbName:       getEnv("DB_NAME", ""),
		CrypteaKey:   getEnv("CRYPTEA_KEY", ""),
	}

	return &config, nil
}

// getEnv is a helper function to read an environment variable or return a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
