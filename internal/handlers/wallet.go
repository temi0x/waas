package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"waas/internal/database"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/joho/godotenv"

	"waas/api"
	"waas/internal/tools/wallet"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

var platformPIN = os.Getenv("CrypteaKey")

func CreateWallet(w http.ResponseWriter, r *http.Request) {

	var params = api.CreateWalletParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	_ = decoder.Decode(&params, r.URL.Query())

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
	// Encode the response
	// encoder := json.NewEncoder(w)
	// encoder.SetIndent("", "  ")
	// encoder.Encode(response)
}

func SendToken(w http.ResponseWriter, r *http.Request) {
	// godotenv.Load(".env")

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

	amount := strconv.FormatFloat(params.Amount, 'f', -1, 64)

	value, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		log.Error("Invalid amount")
		return
	}
	value = new(big.Int).Mul(value, big.NewInt(1000000000000000000)) // Convert to wei// 1 ETH
	gasLimit := uint64(21000)                                        // in units

	toAddress := common.HexToAddress(params.TargetAddress)

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

	// // Specify the recipient and amount
	// recipient := common.HexToAddress("0xRecipientAddress")
	// amount := big.NewInt(1000000000000000000) // 1 ETH

	// // Send the transaction
	// tx, err := client.Transfer(auth, recipient, amount)
	// if err != nil {
	// 	log.Error("Error sending transaction", err)
	// 	return
	// }

	// log.Info("Transaction sent: ", tx.Hash().Hex())

}
