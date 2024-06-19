package api

import (
	"encoding/json"
	"net/http"
)

type CreateWalletParams struct {
	PIN string `json:"pin"` // user's pin, usually 6 digits
}

type SendTokenParams struct {
	UserAddress   string  `json:"user_address"`   // user's wallet address
	PIN           string  `json:"pin"`            // user's pin, usually 6 digits
	TargetAddress string  `json:"target_address"` // e.g., "0x1234..."
	Amount        float64 `json:"amount"`         // e.g., 0.1, 0.3
	ChainID       int     `json:"chain_id"`       // e.g., 1 for Ethereum, 3 for Ropsten
	TokenName     string  `json:"token_name"`     // e.g., "ETH", "USDT"
}

type SendCustomTokenParams struct {
	UserAddress     string  `json:"user_address"`     // user's wallet address
	TargetAddress   string  `json:"target_address"`   // Target address for the custom token
	ContractAddress string  `json:"contract_address"` // Contract address for the custom token
	Amount          float64 `json:"amount"`           // Amount of custom token to send
	PIN             string  `json:"pin"`              // user's pin, usually 6 digits
	ChainID         int     `json:"chain_id"`         // chain ID for request (1, for Ethereum), (3, for Ropsten), etc
}

type SendTokenResponse struct {
	Success bool   `json:"Success"`         // Indicate whether operation was successful
	TxHash  string `json:"TransactionHash"` // Transaction hash of the token transfer
}

type CreateWaasResponse struct {
	Success    bool   `json:"Success"`        // Indicate whether operation was successful
	Address    string `json:"wallet_address"` // Address of the new wallet
	Ciphertext []byte `json:"ciphertext"`     // Encrypted private key
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

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		WriteError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter, err error) {
		WriteError(w, "An Unexpected Error Occurred", http.StatusInternalServerError)
	}
)
