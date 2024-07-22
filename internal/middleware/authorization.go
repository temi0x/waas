package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"waas/api"
	"waas/internal/database"

	// "fmt"
	// "net/http"
	// "github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("invalid username or token")

func ValidateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var apiKey = r.Header.Get("Authorization")

		if apiKey == "" {
			log.Error("No API Key provided")
			api.WriteError(w, "No API Key provided", http.StatusUnauthorized)
			return
		}
		fmt.Println(apiKey)

		// sample api key = CWA-1Ro-ZXSGT5-WlnidH-is5ApkCbi

		// Deconstruct API Key
		keyParts := strings.Split(apiKey, "-")
		if len(keyParts) < 2 {
			log.Error("Invalid API Key format")
			api.WriteError(w, "Invalid API Key format", http.StatusUnauthorized)
			return
		}

		secondGroup := keyParts[1]
		digitsOnly := ""
		for _, char := range secondGroup {
			if unicode.IsDigit(char) {
				digitsOnly += string(char)
			}
		}
		fmt.Println(digitsOnly)

		digitsOnlyInt, err := strconv.Atoi(digitsOnly)
		if err != nil {
			log.Error("Failed to convert digits to int", err)
			api.WriteError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		hashFromDb, err := database.GetKeyHash(digitsOnlyInt)
		if err != nil {
			log.Error("Failed to get hash from db", err)
			api.WriteError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashFromDb), []byte(apiKey))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				log.Error("API Key mismatch", err)
				api.WriteError(w, "Invalid API Key", http.StatusUnauthorized)
			}
			log.Error("Incorrect API Key", err)
			return
		}
		log.Info("Correct API Key")

		next.ServeHTTP(w, r)
	})
}

// func Authorization(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var tokenString = r.Header.Get("Authorization")
// 		var err error

// 		if tokenString == "" {
// 			api.WriteError(w, "No token provided", http.StatusUnauthorized)
// 			return
// 		}

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return tools.JwtSecret, nil
// 		})

// 		if err != nil {
// 			log.Error("Error parsing token", err)
// 			api.WriteError(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			fmt.Println(claims)
// 			ClaimsEmail := claims["email"].(string)

// 			// check if email is in db
// 			var emailExists bool
// 			var err = tools.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM admins WHERE email = ?)", ClaimsEmail).Scan(&emailExists)
// 			if err != nil {
// 				log.Error("Error getting user from database", err)
// 				api.WriteError(w, "Invalid token", http.StatusUnauthorized)
// 				return
// 			}
// 		} else {
// 			log.Error("Error parsing token", err)
// 			api.WriteError(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
