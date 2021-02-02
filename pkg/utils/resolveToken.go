package utils

import (
	"context"
	"errors"
	"strings"
)

type key int

const (
	keyPrincipalID key = iota
)

//GetPrincipalID returns a unique context value name for each hit
func GetPrincipalID() interface{} {
	return keyPrincipalID
}

// GetContextInjected helps get the context data with the specific key given to that context value
func GetContextInjected(ctx context.Context) (string, error) {
	var headerExtracted = ctx.Value(keyPrincipalID)
	if strings.TrimSpace(headerExtracted.(string)) == "" {
		return "", errors.New("AuthorizationFailed")
	}
	return headerExtracted.(string), nil
}
