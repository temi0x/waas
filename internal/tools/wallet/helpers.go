package wallet

import (
	"math/big"
	"strconv"

	etherParams "github.com/ethereum/go-ethereum/params"
	log "github.com/sirupsen/logrus"
)

func ConvertToWei(amount float64) *big.Int {
	// Convert the amount to a string
	amountStr := strconv.FormatFloat(amount, 'f', -1, 64)

	// Convert the string to a big.Int
	amountBigFloat, ok := new(big.Float).SetString(amountStr)
	if !ok {
		log.Fatal("Invalid amount")
	}

	// Multiply the amount by 1 Ether (in Wei) to convert it to Wei
	wei := new(big.Float).Mul(amountBigFloat, new(big.Float).SetFloat64(etherParams.Ether))

	// Convert the big.Float to a big.Int
	weiInt := new(big.Int)
	wei.Int(weiInt)

	return weiInt
}

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
	default:
		return 8453
	}
}

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
	default:
		return "https://basescan.org/tx/" + txHash
	}
}

func GetDecimalPlaces(token string) int {
	switch token {
	case "ETH":
		return 18
	case "BSC":
		return 18
	case "MATIC":
		return 18
	case "FANTOM":
		return 18
	case "XDAI":
		return 18
	case "AVALANCHE":
		return 18
	case "ARBITRUM":
		return 18
	case "OPTIMISM":
		return 18
	case "CELO":
		return 18
	case "MOONBEAM":
		return 18
	case "HARMONY":
		return 18
	case "BASE":
		return 18
	case "USDT":
		return 6
	case "USDC":
		return 6
	case "DAI":
		return 18
	default:
		return 18
	}
}
