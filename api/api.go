package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CreateWalletParams struct {
	PIN string `json:"pin"` // user's pin, usually 6 digits
}

type SendTokenParams struct {
	UserAddress   string  `json:"userAddress"`   // user's wallet address
	PIN           string  `json:"pin"`           // user's pin, usually 6 digits
	TargetAddress string  `json:"targetAddress"` // e.g., "0x1234..."
	Amount        float64 `json:"amount"`        // e.g., 0.1, 0.3
	Chain         string  `json:"chain"`         // e.g., 1 for Ethereum, 3 for Ropsten
	TokenName     string  `json:"tokenName"`     // e.g., "ETH", "USDT"
}

type SendCustomTokenParams struct {
	UserAddress   string  `json:"userAddress"`   // user's wallet address
	TargetAddress string  `json:"targetAddress"` // Target address for the custom token
	TokenName     string  `json:"tokenName"`     // Contract address for the custom token
	Amount        float64 `json:"amount"`        // Amount of custom token to send
	PIN           string  `json:"pin"`           // user's pin, usually 6 digits
	Chain         string  `json:"chain"`         // blockchain to be used for request (ETH, BASE, ARB, etc)
}

type SendTokenResponse struct {
	Success bool   `json:"Success"`         // Indicate whether operation was successful
	TxHash  string `json:"TransactionHash"` // Transaction hash of the token transfer
}

type CreateWaasResponse struct {
	Success    bool   `json:"Success"`       // Indicate whether operation was successful
	Address    string `json:"walletAddress"` // Address of the new wallet
	Ciphertext []byte `json:"ciphertext"`    // Encrypted private key
}

type SwapParams struct {
	UserAddress string  `json:"userAddress"` // user's wallet address
	TargetToken string  `json:"targetToken"` // Target token of the swap operation
	Amount      float64 `json:"amount"`      // amount to be received (in USDT)
	PIN         string  `json:"pin"`         // user's PIN, usually 6 digits
	Chain       string  `json:"chain"`       // blockchain to process transaction on
}

// Error response
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

func HandleError(w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "transfer amount exceeds balance") {
		// Return specific output for this error message
		message := "Insufficient token balance"
		WriteError(w, message, http.StatusUnprocessableEntity)
		return
	}

	if strings.Contains(err.Error(), "no contract code at given address") {
		// Return specific output for this error message
		message := "Incorrect contract address"
		WriteError(w, message, http.StatusBadRequest)
		return
	}

	if strings.Contains(err.Error(), "insufficient funds for gas") {
		// Return specific output for this error message
		message := "Error sending transaction: insufficient funds for gas"
		WriteError(w, message, http.StatusUnprocessableEntity)
		return
	}

	// Handle other error cases here
	// ...
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		HandleError(w, err)
	}
	InternalErrorHandler = func(w http.ResponseWriter, err error) {
		WriteError(w, "An Unexpected Error Occurred", http.StatusInternalServerError)
	}
)
