package pt_client_sdk

import (
	"fmt"
	"net/http"

	"github.com/threatwinds/go-sdk/utils"
)

const (
	AuthAPIURL = "https://inference.threatwinds.com/api/auth/v2/keypair"
)

// ValidateCredentials checks if the provided API credentials are valid
func ValidateCredentials(creds Credentials) error {
	headers := map[string]string{
		"accept":     "application/json",
		"api-key":    creds.APIKey,
		"api-secret": creds.APISecret,
	}

	_, statusCode, err := utils.DoReq[map[string]any](AuthAPIURL, nil, "GET", headers)
	if err != nil {
		return fmt.Errorf("failed to validate credentials: %v", err)
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid credentials (status %d)", statusCode)
	}

	return nil
}
