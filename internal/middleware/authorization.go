package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"waas/api"
	"waas/internal/database"

	// "fmt"
	// "net/http"
	// "github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
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

		// Validate API Key
		isValid, _ := database.ValidateAPIKey(apiKey)
		if !isValid {
			log.Error("Invalid API Keyy")
			api.WriteError(w, "No token provided", http.StatusUnauthorized)
			return
		}
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
