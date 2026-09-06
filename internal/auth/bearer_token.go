package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")

	if bearerToken == "" {
		return "", fmt.Errorf("error getting authorization header")
	}

	tokenString, ok := strings.CutPrefix(strings.TrimSpace(bearerToken), "Bearer ")

	if !ok {
		return "", fmt.Errorf("malformed authorization header")
	}

	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return "", fmt.Errorf("malformed authorization header: empty token")
	}

	return tokenString, nil
}
