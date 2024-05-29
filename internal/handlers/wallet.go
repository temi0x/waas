package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
