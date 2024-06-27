package wallet

import (
	"waas/config"

	log "github.com/sirupsen/logrus"
)

func GetRPC(chainID int) string {
	var rpc string

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load configuration: " + err.Error())
		return ""
	}
	RPC_ETHEREUM := cfg.RPC_ETHEREUM
	RPC_SEPOLIA := cfg.RPC_SEPOLIA

	RPC_POLYGON := cfg.RPC_POLYGON
	RPC_ARBITRUM := cfg.RPC_ARBITRUM
	RPC_OPTIMISM := cfg.RPC_OPTIMISM
	RPC_POLYGON_ZKEVM := cfg.RPC_POLYGON_ZKEVM
	RPC_BASE := cfg.RPC_BASE
	RPC_ASTAR := cfg.RPC_ASTAR
	RPC_ZKSYNC := cfg.RPC_ZKSYNC
	RPC_ZORA := cfg.RPC_ZORA
	RPC_FRAX := cfg.RPC_FRAX
	RPC_ZETACHAIN := cfg.RPC_ZETACHAIN
	RPC_BLAST := cfg.RPC_BLAST

	switch chainID {
	case 1:
		rpc = RPC_ETHEREUM
	case 137:
		rpc = RPC_POLYGON
	case 42161:
		rpc = RPC_ARBITRUM
	case 10:
		rpc = RPC_OPTIMISM
	case 1101:
		rpc = RPC_POLYGON_ZKEVM
	case 8453:
		rpc = RPC_BASE
	case 592:
		rpc = RPC_ASTAR
	case 324:
		rpc = RPC_ZKSYNC
	case 7777777:
		rpc = RPC_ZORA
	case 252:
		rpc = RPC_FRAX
	case 7000:
		rpc = RPC_ZETACHAIN
	case 81457:
		rpc = RPC_BLAST
	case 5:
		rpc = "https://goerli.infura.io/v3/your_infura_project_id"
	case 11155111:
		rpc = RPC_SEPOLIA
	default:
		rpc = "http://localhost:8545"
	}
	return rpc
}
