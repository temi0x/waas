package wallet

import (
	"context"
	"math"
	"strings"

	// "crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"waas/api"
	"waas/config"
	"waas/internal/database"
)

func SendToken(request *api.SendCustomTokenParams) (txhash string, err error) {
	var privateKey []byte
	var RPC_URL string

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %v", err)
	}
	platformPIN := cfg.CrypteaKey

	ChainID := GetChainID(request.Chain)

	RPC_URL = GetRPC(ChainID)

	if RPC_URL == "" {
		return "", fmt.Errorf("unsupported chain: %v", request.Chain)
	}

	if !common.IsHexAddress(request.TargetAddress) {
		return "", fmt.Errorf("invalid target address")
	}

	eNonce, ciphertext, err := database.GetWalletDetails(request.UserAddress)
	if err != nil {
		return "", fmt.Errorf("error getting wallet details: %v", err)
	}

	derivedKey, err := Decrypt(eNonce, ciphertext, request.PIN, platformPIN)
	if err != nil {
		return "", fmt.Errorf("error decrypting user wallet: %v", err)
	}
	derivedKeyWithPrefix := fmt.Sprintf("0x%s", derivedKey)

	privateKey, err = hexutil.Decode(derivedKeyWithPrefix)
	if err != nil {
		return "", fmt.Errorf("error decoding private key: %v", err)
	}

	privKey, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("error converting to ECDSA private key: %v", err)
	}

	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", fmt.Errorf("error connecting to Ethereum client: %v", err)
	}

	chainID := big.NewInt(int64(ChainID))
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return "", fmt.Errorf("error creating transactor: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return "", fmt.Errorf("error getting nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("error getting gas price: %v", err)
	}

	amount := ConvertToWei(request.Amount)
	value := new(big.Int).Set(amount)
	gasLimit := uint64(21000) // in units

	toAddress := common.HexToAddress(request.TargetAddress)
	if toAddress == (common.Address{}) {
		return "", fmt.Errorf("invalid target address")
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %v", err)
	}

	txHash := GetBlockExplorerURL(ChainID, signedTx.Hash().Hex())

	return txHash, nil
}

// SendTokens performs the custom token transfer based on the provided request details
func SendTokens(request *api.SendCustomTokenParams, decimals int, contractAddress string) (txhash string, err error) {
	ChainID := GetChainID(request.Chain)

	rpcURL := GetRPC(ChainID)
	if rpcURL == "" {
		return "", fmt.Errorf("unsupported chain ID: %v", request.Chain)
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	chainID := big.NewInt(int64(ChainID)) // Chain ID for Ethereum mainnet

	// Decrypt the private key
	eNonce, ciphertext, err := database.GetWalletDetails(request.UserAddress)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve wallet details: %w", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load configuration: " + err.Error())
		return "", err
	}
	platformPIN := cfg.CrypteaKey
	decryptedKey, err := Decrypt(eNonce, ciphertext, request.PIN, platformPIN)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(decryptedKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Prepare the transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %w", err)
	}

	// Calculate the gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	// Load the default ERC-20 ABI
	parsedABI, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse ERC-20 ABI: %w", err)
	}

	// Create an instance of the token contract
	tokenAddress := common.HexToAddress(contractAddress)
	token := bind.NewBoundContract(tokenAddress, parsedABI, client, client, client)

	// Convert the amount to the smallest unit by multiplying by 10^decimals
	amountInSmallestUnit := new(big.Float).Mul(new(big.Float).SetFloat64(request.Amount), new(big.Float).SetFloat64(math.Pow10(decimals)))

	// Convert the amount to *big.Int
	amountBigInt := new(big.Int)
	amountInSmallestUnit.Int(amountBigInt) // Convert the float to *big.Int without losing precision
	// Use `amountBigInt` in the call
	// Prepare the transaction
	tx, err := token.Transact(auth, "transfer", common.HexToAddress(request.TargetAddress), amountBigInt)
	if err != nil {
		return "", fmt.Errorf("failed to send token: %w", err)
	}

	log.Printf("Token sent: %s", tx.Hash().Hex())
	txHash := GetBlockExplorerURL(ChainID, tx.Hash().Hex())
	return txHash, nil
}
