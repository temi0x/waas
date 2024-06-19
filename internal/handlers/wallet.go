package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	database "waas/internal/database"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	etherParams "github.com/ethereum/go-ethereum/params"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	"waas/api"
	"waas/config"
	"waas/internal/tools/wallet"
)

func CreateWallet(w http.ResponseWriter, r *http.Request) {

	var params = api.CreateWalletParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	_ = decoder.Decode(&params, r.URL.Query())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load configuration: " + err.Error())
		api.InternalErrorHandler(w, err)
		return
	}
	platformPIN := cfg.CrypteaKey

	// Create a new wallet
	pKey, address := wallet.CreateWallet()
	var walletAddress = address

	//Encrypt the wallet
	nonce, ciphertext, key, err := wallet.Encrypt(pKey, params.PIN, platformPIN)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	// Store the wallet in the database
	err = database.StoreWalletDetails(nonce, address, key, ciphertext)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w, err)
		return
	}

	// Create a new response
	var response = api.CreateWaasResponse{
		Success:    true,
		Address:    walletAddress,
		Ciphertext: ciphertext,
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

func SendToken(w http.ResponseWriter, r *http.Request) {

	var params = api.SendTokenParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var privateKey []byte
	var RPC_URL string
	var RPC_ETHEREUM, RPC_RINKEBY, RPC_ROPSTEN, RPC_SEPOLIA string
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
	RPC_ETHEREUM = cfg.RPC_ETHEREUM
	RPC_RINKEBY = cfg.RPC_RINKEBY
	RPC_ROPSTEN = cfg.RPC_ROPSTEN
	RPC_SEPOLIA = cfg.RPC_SEPOLIA

	// set RPC URL based on chain ID
	switch params.ChainID {
	case 1:
		RPC_URL = RPC_ETHEREUM
	case 3:
		RPC_URL = RPC_ROPSTEN
	case 4:
		RPC_URL = RPC_RINKEBY
	case 11155111:
		RPC_URL = RPC_SEPOLIA
	default:
		api.WriteError(w, "Unsupported chain ID", http.StatusBadRequest)
		return
	}

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
	derivedKey, err := wallet.Decrypt(eNonce, ciphertext, params.PIN, platformPIN)
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

	amount := convertToWei(params.Amount)

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

func SendCustomToken(w http.ResponseWriter, r *http.Request) {

	var params = api.SendCustomTokenParams{}
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

	// set RPC URL based on chain ID
	switch params.ChainID {
	case 1:
		RPC_URL = os.Getenv("RPC_ETHEREUM")
	case 3:
		RPC_URL = os.Getenv("RPC_ROPSTEN")
	case 4:
		RPC_URL = os.Getenv("RPC_RINKEBY")
	case 11155111:
		RPC_URL = os.Getenv("RPC_SEPOLIA")
	default:
		api.WriteError(w, "Unsupported chain ID", http.StatusBadRequest)
		return
	}

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
	derivedKey, err := wallet.Decrypt(eNonce, ciphertext, params.PIN, platformPIN)
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

	amount := convertToWei(params.Amount)

	value := new(big.Int).Set(amount)

	gasLimit := uint64(21000) // in units

	toAddress := common.HexToAddress(params.TargetAddress)
	if toAddress == (common.Address{}) {
		log.Error("Invalid target address")
		return
	}

	// Create and sign the transaction
	tx := types.NewTransaction(nonce, common.HexToAddress(params.ContractAddress), value, gasLimit, gasPrice, createTransferData(toAddress, amount))

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(int64(params.ChainID))), privKey)
	if err != nil {
		http.Error(w, "Error signing transaction", http.StatusInternalServerError)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		http.Error(w, "Error sending transaction", http.StatusInternalServerError)
		return
	}

	log.Info("Transaction sent: ", signedTx.Hash().Hex())
}

func createTransferData(to common.Address, amount *big.Int) []byte {
	// ABI encoded function signature and parameters for ERC-20 transfer function
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data
}

func convertToWei(amount float64) *big.Int {
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
