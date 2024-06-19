package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RPC_ETHEREUM  string
	RPC_OPTIMISM  string
	RPC_POLYGON   string
	RPC_BSC       string
	RPC_FANTOM    string
	RPC_XDAI      string
	RPC_ARBITRUM  string
	RPC_MOONBEAM  string
	RPC_AVALANCHE string

	RPC_SOLANA string
	RPC_COSMOS string

	RPC_SEPOLIA string
	RPC_RINKEBY string
	RPC_ROPSTEN string
	SqlUsername string
	SqlPassword string
	SqlUrl      string
	DbName      string
	CrypteaKey  string
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
		RPC_ETHEREUM:  getEnv("RPC_ETHEREUM", ""),
		RPC_SEPOLIA:   getEnv("RPC_SEPOLIA", ""),
		RPC_RINKEBY:   getEnv("RPC_RINKEBY", ""),
		RPC_ROPSTEN:   getEnv("RPC_ROPSTEN", ""),
		RPC_OPTIMISM:  getEnv("RPC_OPTIMISM", ""),
		RPC_POLYGON:   getEnv("RPC_POLYGON", ""),
		RPC_BSC:       getEnv("RPC_BSC", ""),
		RPC_FANTOM:    getEnv("RPC_FANTOM", ""),
		RPC_XDAI:      getEnv("RPC_XDAI", ""),
		RPC_ARBITRUM:  getEnv("RPC_ARBITRUM", ""),
		RPC_MOONBEAM:  getEnv("RPC_MOONBEAM", ""),
		RPC_AVALANCHE: getEnv("RPC_AVALANCHE", ""),

		RPC_SOLANA: getEnv("RPC_SOLANA", ""),
		RPC_COSMOS: getEnv("RPC_COSMOS", ""),

		SqlUsername: getEnv("SQL_USERNAME", ""),
		SqlPassword: getEnv("SQL_PASSWORD", ""),
		SqlUrl:      getEnv("SQL_URL", ""),
		DbName:      getEnv("DB_NAME", ""),
		CrypteaKey:  getEnv("CRYPTEA_KEY", ""),
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
