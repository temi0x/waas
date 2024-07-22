package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"

	"waas/api"
	"waas/config"
	database "waas/internal/database"
	"waas/internal/tools/analytics"
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
	//Decode the request
	var request api.SendCustomTokenParams
	var err error
	if err = schema.NewDecoder().Decode(&request, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start := time.Now()

	var response api.SendTokenResponse

	//check if token is native or custom
	if request.TokenName == "ETH" {
		// Perform the token transfer
		txHash, err := wallet.SendToken(&request)
		if err != nil {
			recordMetrics("failed", time.Since(start).Seconds())

			analytics.StoreTransaction(analytics.TransactionLog{
				WalletAddress: request.UserAddress,
				TargetAddress: request.TargetAddress,
				TokenType:     "ERC20",
				Amount:        fmt.Sprintf("%f", request.Amount),
				Status:        "failed",
				ErrorMessage:  fmt.Sprintf("Failed to send token: %v", err),
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			})
			log.Printf("Failed to send token: %v", err)
			http.Error(w, "Failed to send token", http.StatusInternalServerError)
			return
		} else {
			recordMetrics("success", time.Since(start).Seconds())
			analytics.StoreTransaction(analytics.TransactionLog{
				WalletAddress: request.UserAddress,
				TargetAddress: request.TargetAddress,
				TxnID:         txHash,
				TokenType:     "ERC20",
				Amount:        fmt.Sprintf("%f", request.Amount),
				Status:        "success",
				ErrorMessage:  "",
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			})
		}
		response.Success = err == nil
		response.TxHash = txHash
	} else {
		//check and validate necessary parameters
		decimals := wallet.GetDecimalPlaces(request.TokenName)

		// validate contract address with chain
		tokenContractAddress := wallet.GetContractAddress(request.Chain, request.TokenName)
		log.Println("Token contract address: ", tokenContractAddress)

		// Perform the token transfer
		txHash, err := wallet.SendTokens(&request, decimals, tokenContractAddress)

		if err != nil {
			recordMetrics("failed", time.Since(start).Seconds())
			analytics.StoreTransaction(analytics.TransactionLog{
				WalletAddress: request.UserAddress,
				TargetAddress: request.TargetAddress,
				TokenType:     "ERC20",
				Amount:        fmt.Sprintf("%f", request.Amount),
				Status:        "failed",
				ErrorMessage:  fmt.Sprintf("Failed to send token: %v", err),
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			})
			log.Printf("Failed to send token: %v", err)
			http.Error(w, "Failed to send token", http.StatusInternalServerError)
			return

		} else {
			recordMetrics("success", time.Since(start).Seconds())
			analytics.StoreTransaction(analytics.TransactionLog{
				WalletAddress: request.UserAddress,
				TargetAddress: request.TargetAddress,
				TxnID:         txHash,
				TokenType:     "ERC20",
				Amount:        fmt.Sprintf("%f", request.Amount),
				Status:        "success",
				ErrorMessage:  "",
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			})
		}
		response.Success = err == nil
		response.TxHash = txHash
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

// // SendToken handles the HTTP request to transfer a token.
// func SendTokens(w http.ResponseWriter, r *http.Request) {
// 	// Decode the request
// 	var request api.SendCustomTokenParams
// 	if err := schema.NewDecoder().Decode(&request, r.URL.Query()); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	start := time.Now()

// 	var response api.SendTokenResponse

// 	//check and validate necessary parameters
// 	decimals := wallet.GetDecimalPlaces(request.TokenName)

// 	// validate contract address with chain
// 	tokenContractAddress := wallet.GetContractAddress(request.Chain, request.TokenName)
// 	log.Println("Token contract address: ", tokenContractAddress)

// 	// Perform the token transfer
// 	txHash, err := wallet.SendTokens(&request, decimals, tokenContractAddress) // Adjusted to pass a pointer and handle both return values
// 	if err != nil {
// 		recordMetrics("failed", time.Since(start).Seconds())

// 		analytics.StoreTransaction(analytics.TransactionLog{
// 			WalletAddress: request.UserAddress,
// 			TargetAddress: request.TargetAddress,
// 			TokenType:     "ERC20",
// 			Amount:        fmt.Sprintf("%f", request.Amount),
// 			Status:        "failed",
// 			ErrorMessage:  fmt.Sprintf("Failed to send token: %v", err),
// 			Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
// 		})
// 		response.Success = false
// 		response.TxHash = txHash

// 		log.Printf("Failed to send token: %v", err)
// 		// http.Error(w, "Failed to send token", http.StatusInternalServerError)
// 		return
// 	} else {
// 		recordMetrics("success", time.Since(start).Seconds())

// 		analytics.StoreTransaction(analytics.TransactionLog{
// 			WalletAddress: request.UserAddress,
// 			TargetAddress: request.TargetAddress,
// 			TokenType:     "ERC20",
// 			Amount:        fmt.Sprintf("%f", request.Amount),
// 			Status:        "success",
// 			ErrorMessage:  "",
// 			Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
// 		})
// 		response.Success = true
// 		response.TxHash = txHash
// 	}

// 	// Set the header to application/json
// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		log.Error(err)
// 		api.InternalErrorHandler(w, err)
// 		return
// 	}

// 	if !response.Success {
// 		http.Error(w, "Failed to send token", http.StatusInternalServerError)
// 		return
// 	}
// }

func recordMetrics(status string, duration float64) {
	totalTransactions.WithLabelValues(status).Inc()
	transactionDuration.WithLabelValues(status).Observe(duration)
}

func recordWalletCreationMetrics(status string, duration float64, walletAddress string, business string) {
	walletCreation.WithLabelValues(status, walletAddress, business).Observe(duration)
}
