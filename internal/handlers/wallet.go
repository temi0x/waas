package handlers

import (
	"net/http"
	"encoding/json"

	"waas/api"
	"waas/internal/tools/wallet"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	// Set the header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Create a new wallet
	privateKey, address := wallet.CreateWallet()

	// Create a new response
	response := api.Response{
		Success: true,
		address: address,
		// Data: map[string]string{
		// 	"private_key": privateKey,
		// 	"address": address,
		// },
	}

	// Encode the response
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(response)
}

func EncryptPkey (privateKey string, pin string) {
	// Encrypt the private key
	encryptedKey := wallet.EncryptPkey(privateKey, pin)

	// Create a new response
	response := api.Response{
		Success: true,
		Data: map[string]string{
			"encrypted_key": encryptedKey,
		},
	}

	// Encode the response
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(response)
}