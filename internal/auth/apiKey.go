package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	// Check if header exists and starts with "ApiKey "
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("Unauthorized")
	}

	return strings.TrimPrefix(authHeader, "ApiKey "), nil
}