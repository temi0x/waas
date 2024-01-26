package middleware

import (
	"errors"
	// "fmt"
	// "net/http"

	// "github.com/go-chi/chi"
	// log "github.com/sirupsen/logrus"
)

var UnAuthorizedError = errors.New("Invalid username or token")

