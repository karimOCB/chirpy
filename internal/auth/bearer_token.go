package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")

	if bearerToken == "" {
		return "", fmt.Errorf("error getting bearer token")
	}

	tokenString, ok := strings.CutPrefix(bearerToken, "Bearer ")
	if !ok {
		return 
	}  
}