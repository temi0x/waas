package api

import (
	"encoding/json"
	"net/http"
)

type CreateWalletParams struct {
	PIN string
}

type CreateWaasResponse struct {
	Success    bool
	Address    string
	Ciphertext []byte
}

// Error response
type Error struct {
	Code    int
	Message string
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
