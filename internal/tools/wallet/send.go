package wallet

import (
	"context"
	"math"
	"strings"

	// "crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"

	"waas/api"
	"waas/config"
	"waas/internal/database"
)

func SendToken(w http.ResponseWriter, r *http.Request) {

	var params = api.SendTokenParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var privateKey []byte
	var RPC_URL string
	var err error

	err = decoder.Decode(&params, r.URL.Query())
	if err != nil {
		api.InternalErrorHandler(w, err)
		log.Error("Error decoding request params", err)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load configuration: " + err.Error())
		api.InternalErrorHandler(w, err)
		return
	}
	platformPIN := cfg.CrypteaKey
	RPC_URL = GetRPC(params.ChainID)

	log.Info("RPC_URL: ", RPC_URL)

	if RPC_URL == "" {
		api.WriteError(w, "RPC URL is not set", http.StatusInternalServerError)
		return
	}

	// check if targetaddress is EVM-compatible
	if !common.IsHexAddress(params.TargetAddress) {
		api.WriteError(w, "Invalid target address", http.StatusBadRequest)
		return
	}

	// get nonce & ciphertext from database
	eNonce, ciphertext, err := database.GetWalletDetails(params.UserAddress)
	if err != nil {
		api.InternalErrorHandler(w, err)
		log.Error("Error getting wallet details", err)
		return
	}

	// Decrypt the wallet
	derivedKey, err := Decrypt(eNonce, ciphertext, params.PIN, platformPIN)
	log.Info("derivedKey: ", derivedKey)
	if err != nil {
		api.InternalErrorHandler(w, err)
		log.Error("Error decrypting user wallet", err)
		return
	}
	derivedKeyWithPrefix := fmt.Sprintf("0x%s", derivedKey)

	// privateKeyBytes := []byte(derivedKey)

	privateKey, err = hexutil.Decode(derivedKeyWithPrefix)
	if err != nil {
		log.Error("Error decoding private key", err)
		return
	}
	log.Info("private key: ", privateKey)

	privKey, err := crypto.ToECDSA([]byte(privateKey))
	log.Info("private key ecdsa: ", privKey)
	if err != nil {
		log.Error("Error converting to ECDSA private key", err)
		return
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Error("Error connecting to Ethereum client ", err)
		return
	}

	// Create a new authenticated session
	chainID := big.NewInt(int64(params.ChainID)) // 1 is the chain ID for Ethereum mainnet
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		log.Error("Error creating transactor", err)
		return
	}

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Error("Error getting nonce", err)
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("Error getting gas price", err)
		return
	}

	amount := ConvertToWei(params.Amount)

	value := new(big.Int).Set(amount)

	gasLimit := uint64(21000) // in units

	toAddress := common.HexToAddress(params.TargetAddress)
	if toAddress == (common.Address{}) {
		log.Error("Invalid target address")
		return
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		log.Error("Error signing transaction", err)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Error("Error sending transaction", err)
		return
	}

	log.Info("Transaction sent: ", signedTx.Hash().Hex())

	txHash := signedTx.Hash().Hex()

	// Create a new response
	var response = api.SendTokenResponse{
		Success: true,
		TxHash:  txHash,
	}

	// Set the header to application/json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w, err)
		return
	}

}

const erc20ABI = `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

// SendTokens performs the custom token transfer based on the provided request details
func SendTokens(request *api.SendCustomTokenParams) (txhash string, err error) {
	// Determine the RPC URL based on the Chain ID
	rpcURL := GetRPC(request.ChainID)
	if rpcURL == "" {
		return "", fmt.Errorf("unsupported chain ID: %d", request.ChainID)
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	chainID := big.NewInt(int64(request.ChainID)) // Chain ID for Ethereum mainnet

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
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse ERC-20 ABI: %w", err)
	}

	// Create an instance of the token contract
	tokenAddress := common.HexToAddress(request.ContractAddress)
	token := bind.NewBoundContract(tokenAddress, parsedABI, client, client, client)

	const tokenDecimals = 18

	// Convert the amount to the smallest unit by multiplying by 10^decimals
	amountInSmallestUnit := new(big.Float).Mul(new(big.Float).SetFloat64(request.Amount), new(big.Float).SetFloat64(math.Pow10(tokenDecimals)))

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
	txHash := tx.Hash().Hex()
	return txHash, nil
}
