package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"
	"waas/config"

	etherParams "github.com/ethereum/go-ethereum/params"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// ConvertToWei converts an amount in Ether to Wei
func ConvertToWei(amount float64) *big.Int {
	amountStr := strconv.FormatFloat(amount, 'f', -1, 64)
	amountBigFloat, ok := new(big.Float).SetString(amountStr)
	if !ok {
		log.Fatal("Invalid amount")
	}

	wei := new(big.Float).Mul(amountBigFloat, new(big.Float).SetFloat64(etherParams.Ether))
	weiInt := new(big.Int)
	wei.Int(weiInt)

	return weiInt
}

// GenerateTransactionID generates a unique transaction ID
func GenerateTransactionID() string {
	id := uuid.New()
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d-%s", timestamp, id.String())
}

// GetChainID returns the chain ID for a given chain name
func GetChainID(chain string) int {
	switch chain {
	case "ETH":
		return 1
	case "BSC":
		return 56
	case "MATIC":
		return 137
	case "FANTOM":
		return 250
	case "XDAI":
		return 100
	case "AVALANCHE":
		return 43114
	case "ARBITRUM":
		return 42161
	case "OPTIMISM":
		return 10
	case "CELO":
		return 42220
	case "MOONBEAM":
		return 1287
	case "HARMONY":
		return 1666600000
	case "BASE":
		return 8453
	case "FVM_Testnet":
		return 314159
	case "SEPOLIA":
		return 11155111
	default:
		return 8453
	}
}

// GetBlockExplorerURL returns the block explorer URL for a given chain ID and transaction hash
func GetBlockExplorerURL(chainID int, txHash string) string {
	switch chainID {
	case 1:
		return "https://etherscan.io/tx/" + txHash
	case 56:
		return "https://bscscan.com/tx/" + txHash
	case 137:
		return "https://polygonscan.com/tx/" + txHash
	case 250:
		return "https://ftmscan.com/tx/" + txHash
	case 100:
		return "https://blockscout.com/xdai/mainnet/tx/" + txHash
	case 43114:
		return "https://cchain.explorer.avax.network/tx/" + txHash
	case 42161:
		return "https://arbiscan.io/tx/" + txHash
	case 10:
		return "https://optimistic.etherscan.io/tx/" + txHash
	case 42220:
		return "https://explorer.celo.org/tx/" + txHash
	case 1287:
		return "https://moonbeam-explorer.netlify.app/tx/" + txHash
	case 1666600000:
		return "https://explorer.harmony.one/tx/" + txHash
	case 8453:
		return "https://basescan.org/tx/" + txHash
	case 314159:
		return "https://calibration.filscan.io/en/message/" + txHash
	case 11155111:
		return "https://sepolia.etherscan.io/tx/" + txHash
	default:
		return "https://basescan.org/tx/" + txHash
	}
}

// GetDecimalPlaces returns the number of decimal places for a given token
func GetDecimalPlaces(token string) int {
	switch token {
	case "ETH", "BSC", "MATIC", "FANTOM", "XDAI", "AVALANCHE", "ARBITRUM", "OPTIMISM", "CELO", "MOONBEAM", "HARMONY", "BASE", "DAI":
		return 18
	case "USDT", "USDC":
		return 6
	default:
		return 18
	}
}

// GetTokenRate fetches the current rate of the specified token from CoinGecko
func GetTokenRate(tokenName string) float64 {
	configs, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed to load configuration: %v", err)
		return 0
	}
	coinGeckoKey := configs.CoinGeckoKey

	token := mapTokenName(tokenName)
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return 0
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", coinGeckoKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to execute request: %v", err)
		return 0
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return 0
	}

	var response map[string]map[string]float64
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("failed to unmarshal response: %v", err)
		return 0
	}

	rate, ok := response[token]["usd"]
	if !ok {
		log.Printf("failed to get token rate for %s", tokenName)
		return 0
	}

	return rate
}

// mapTokenName maps common token names to their CoinGecko IDs
func mapTokenName(tokenName string) string {
	switch tokenName {
	case "ETH":
		return "ethereum"
	case "USDC":
		return "usd-coin"
	case "USDT":
		return "tether"
	case "DAI":
		return "dai"
	default:
		return tokenName
	}
}
