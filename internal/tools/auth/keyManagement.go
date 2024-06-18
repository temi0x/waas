package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateAPIKey(userID string) string {
	// Get current time in microseconds
	currentTime := time.Now().UnixNano() / int64(time.Microsecond)

	// Create a string with the user ID, current time, and "Cryptea"
	data := fmt.Sprintf("%s%dCryptea", userID, currentTime)

	// Create a new SHA256 hash
	hash := sha256.New()
	hash.Write([]byte(data))

	// Convert the hash to a hexadecimal string
	apiKey := hex.EncodeToString(hash.Sum(nil))

	return apiKey
}
