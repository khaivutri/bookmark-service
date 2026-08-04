package infrastructure

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func loadKeyFromEnv(envName string) ([]byte, error) {
	value := strings.TrimSpace(os.Getenv(envName))
	if value == "" {
		return nil, fmt.Errorf("missing env %s", envName)
	}

	// Base64 is convenient for .env files; raw PEM is also accepted.
	if decoded, err := base64.StdEncoding.DecodeString(value); err == nil {
		return decoded, nil
	}
	return []byte(value), nil
}