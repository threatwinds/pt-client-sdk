package twpt_client_sdk

import (
	"net/http"
	"sync"
)

type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	Credentials *Credentials
}

var (
	instance *Client
	once     sync.Once
)

func NewTWPTClient(host string, creds Credentials) (*Client, error) {
	once.Do(func() {
		instance = &Client{
			BaseURL:     host,
			HTTPClient:  &http.Client{},
			Credentials: &creds,
		}
	})

	return instance, nil
}

func GetTWPTClient() *Client {
	return instance
}
