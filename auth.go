package main

import (
	"errors"
	"net/http"
	"strings"
)

const NoAuthHeaderIncluded = "no Auth header included in request"
const Apikey = "ApiKey"

func GetToken(headers http.Header, tokenName string) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New(NoAuthHeaderIncluded)
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != tokenName {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
