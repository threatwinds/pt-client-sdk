package twpt_client_sdk

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/threatwinds/go-sdk/utils"
)

const (
	AuthAPIURL = "https://inference.threatwinds.com/api/auth/v2/keypair"
)

type Credentials struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	Credentials *Credentials
}

var (
	instance *Client
	once     sync.Once
)

func NewTWPTClient(host, apiKey, apiSecret string) (*Client, error) {
	once.Do(func() {
		instance = &Client{
			BaseURL:    host,
			HTTPClient: &http.Client{},
			Credentials: &Credentials{
				APIKey:    apiKey,
				APISecret: apiSecret,
			},
		}
	})

	return instance, nil
}

func GetTWPTClient() *Client {
	return instance
}

func (c *Client) ValidateCredentials() error {
	headers := map[string]string{
		"accept":     "application/json",
		"api-key":    c.Credentials.APIKey,
		"api-secret": c.Credentials.APISecret,
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
