package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: ApiKey API_KEY_VALUE
func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("not authorized")
	}

	headerAPIKey := strings.Split(header, " ")
	if len(headerAPIKey) != 2 || headerAPIKey[0] != "ApiKey" {
		return "", errors.New("malformed authorization header value")
	}

	return headerAPIKey[1], nil
}
