package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RPC_ETHEREUM      string
	RPC_POLYGON       string
	RPC_ARBITRUM      string
	RPC_OPTIMISM      string
	RPC_POLYGON_ZKEVM string
	RPC_BASE          string
	RPC_ASTAR         string
	RPC_ZKSYNC        string
	RPC_ZORA          string
	RPC_FRAX          string
	RPC_ZETACHAIN     string
	RPC_BLAST         string

	RPC_BSC       string
	RPC_FANTOM    string
	RPC_XDAI      string
	RPC_MOONBEAM  string
	RPC_AVALANCHE string

	RPC_SOLANA string
	RPC_COSMOS string

	RPC_SEPOLIA              string
	RPC_HOLESKY              string
	RPC_POLYGON_AMOY         string
	RPC_ARB_SEPOLIA          string
	RPC_OP_SEPOLIA           string
	RPC_POLYGONZK_CARDONA    string
	RPC_BASE_SEPOLIA         string
	RPC_ZKSYNC_SEPOLIA       string
	RPC_ZORA_SEPOLIA         string
	RPC_FRAX_SEPOLIA         string
	RPC_ZETACHAIN_SEPOLIA    string
	RPC_BLAST_SEPOLIA        string
	RPC_FILECOIN_CALIBRATION string

	SqlUsername  string
	SqlPassword  string
	SqlUrl       string
	DbName       string
	CrypteaKey   string
	CoinGeckoKey string
}

// LoadConfig reads configuration from .env file and environment variables.
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file, loading config from environment")
	}

	config := Config{
		RPC_ETHEREUM:      getEnv("RPC_ETHEREUM", ""),
		RPC_POLYGON:       getEnv("RPC_POLYGON", ""),
		RPC_ARBITRUM:      getEnv("RPC_ARBITRUM", ""),
		RPC_OPTIMISM:      getEnv("RPC_OPTIMISM", ""),
		RPC_POLYGON_ZKEVM: getEnv("RPC_POLYGON_ZKEVM", ""),
		RPC_BASE:          getEnv("RPC_BASE", ""),
		RPC_ASTAR:         getEnv("RPC_ASTAR", ""),
		RPC_ZKSYNC:        getEnv("RPC_ZKSYNC", ""),
		RPC_ZORA:          getEnv("RPC_ZORA", ""),
		RPC_FRAX:          getEnv("RPC_FRAX", ""),
		RPC_ZETACHAIN:     getEnv("RPC_ZETACHAIN", ""),
		RPC_BLAST:         getEnv("RPC_BLAST", ""),

		RPC_BSC:       getEnv("RPC_BSC", ""),
		RPC_FANTOM:    getEnv("RPC_FANTOM", ""),
		RPC_XDAI:      getEnv("RPC_XDAI", ""),
		RPC_MOONBEAM:  getEnv("RPC_MOONBEAM", ""),
		RPC_AVALANCHE: getEnv("RPC_AVALANCHE", ""),

		RPC_SEPOLIA:              getEnv("RPC_SEPOLIA", ""),
		RPC_HOLESKY:              getEnv("RPC_HOLESKY", ""),
		RPC_POLYGON_AMOY:         getEnv("RPC_POLYGON_AMOY", ""),
		RPC_ARB_SEPOLIA:          getEnv("RPC_ARB_SEPOLIA", ""),
		RPC_OP_SEPOLIA:           getEnv("RPC_OP_SEPOLIA", ""),
		RPC_POLYGONZK_CARDONA:    getEnv("RPC_POLYGONZK_CARDONA", ""),
		RPC_BASE_SEPOLIA:         getEnv("RPC_BASE_SEPOLIA", ""),
		RPC_ZKSYNC_SEPOLIA:       getEnv("RPC_ZKSYNC_SEPOLIA", ""),
		RPC_ZORA_SEPOLIA:         getEnv("RPC_ZORA_SEPOLIA", ""),
		RPC_FRAX_SEPOLIA:         getEnv("RPC_FRAX_SEPOLIA", ""),
		RPC_ZETACHAIN_SEPOLIA:    getEnv("RPC_ZETACHAIN_SEPOLIA", ""),
		RPC_BLAST_SEPOLIA:        getEnv("RPC_BLAST_SEPOLIA", ""),
		RPC_FILECOIN_CALIBRATION: getEnv("RPC_FVM_CALIBRATION", ""),

		RPC_SOLANA: getEnv("RPC_SOLANA", ""),
		RPC_COSMOS: getEnv("RPC_COSMOS", ""),

		SqlUsername:  getEnv("SQL_USERNAME", ""),
		SqlPassword:  getEnv("SQL_PASSWORD", ""),
		SqlUrl:       getEnv("SQL_URL", ""),
		DbName:       getEnv("DB_NAME", ""),
		CrypteaKey:   getEnv("CRYPTEA_KEY", ""),
		CoinGeckoKey: getEnv("GECKO_API_KEY", ""),
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
