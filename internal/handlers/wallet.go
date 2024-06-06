package handlers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"waas/api"
	"waas/internal/tools"
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
	err = tools.StoreWalletDetails(nonce, address, key, ciphertext)
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
	var params = api.SendTokenParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var privateKey []byte
	var err error

	_ = decoder.Decode(&params, r.URL.Query())
	if err != nil {
		api.InternalErrorHandler(w, err)
		log.Error("Error decoding login params", err)
		return
	}

	// check if targetaddress is EVM-compatible
	if !common.IsHexAddress(params.TargetAddress) {
		api.WriteError(w, "Invalid target address", http.StatusBadRequest)
		return
	}

	// get nonce from database
	// Decrypt the wallet
	nonce, ciphertext, err := tools.GetWalletDetails(params.UserAddress)
	// ciphertext, err = wallet.Decrypt(nonce, ciphertext, params.PIN, platformPIN)

	derivedKey, err := wallet.Decrypt(nonce, ciphertext, params.PIN, platformPIN)
	if err != nil {
		api.InternalErrorHandler(w, err)
		log.Error("Error decrypting user wallet", err)
		return
	}
	privateKey = []byte(derivedKey)

	privKey, err := crypto.ToECDSA(privateKey)
	if err != nil {
		log.Error("Error converting to ECDSA private key", err)
		return
	}

	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Error("Error connecting to Ethereum client", err)
		return
	}

	// Create a new authenticated session
	chainID := big.NewInt(1) // 1 is the chain ID for Ethereum mainnet
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		log.Error("Error creating transactor", err)
		return
	}

	// Specify the recipient and amount
	recipient := common.HexToAddress("0xRecipientAddress")
	amount := big.NewInt(1000000000000000000) // 1 ETH

	// Send the transaction
	// tx, err := client.Transfer(auth, recipient, amount)
	// if err != nil {
	// 	log.Error("Error sending transaction", err)
	// 	return
	// }

	log.Info("Transaction sent: ", tx.Hash().Hex())

}
