package pt_client_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/threatwinds/go-sdk/utils"
	"github.com/threatwinds/pt-client-sdk/helpers"
)

// HTTPClient provides HTTP access to the ThreatWinds Pentest API
type HTTPClient struct {
	BaseURL     string
	HTTPClient  *http.Client
	Credentials *Credentials
}

// NewHTTPClient creates a new HTTP client instance
func NewHTTPClient(baseURL string, creds Credentials) *HTTPClient {
	return &HTTPClient{
		BaseURL:     baseURL,
		HTTPClient:  &http.Client{},
		Credentials: &creds,
	}
}

// ListPentests retrieves a paginated list of pentests
func (c *HTTPClient) ListPentests(ctx context.Context, pagination PaginationParams) (*PentestListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/pentests?page=%d&page_size=%d", c.BaseURL, pagination.Page, pagination.PageSize)

	headers := map[string]string{
		"accept":     "application/json",
		"api-key":    c.Credentials.APIKey,
		"api-secret": c.Credentials.APISecret,
	}

	result, statusCode, err := utils.DoReq[PentestListResponse](url, nil, "GET", headers)
	if err != nil {
		return nil, fmt.Errorf("failed to list pentests: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return &result, nil
}

// GetPentest retrieves a single pentest by ID
func (c *HTTPClient) GetPentest(ctx context.Context, pentestID string) (*PentestData, error) {
	url := fmt.Sprintf("%s/api/v1/pentests/%s", c.BaseURL, pentestID)

	headers := map[string]string{
		"accept":     "application/json",
		"api-key":    c.Credentials.APIKey,
		"api-secret": c.Credentials.APISecret,
	}

	result, statusCode, err := utils.DoReq[PentestData](url, nil, "GET", headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get pentest: %w", err)
	}

	if statusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pentest %s not found", pentestID)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return &result, nil
}

// SchedulePentest schedules a new pentest
func (c *HTTPClient) SchedulePentest(ctx context.Context, req *SchedulePentestRequest) (string, error) {
	url := fmt.Sprintf("%s/api/v1/pentests/schedule", c.BaseURL)

	headers := map[string]string{
		"accept":       "application/json",
		"content-type": "application/json",
		"api-key":      c.Credentials.APIKey,
		"api-secret":   c.Credentials.APISecret,
	}

	bodyJson, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	result, statusCode, err := utils.DoReq[SchedulePentestResponse](url, bodyJson, "POST", headers)
	if err != nil {
		return "", fmt.Errorf("failed to schedule pentest: %w", err)
	}

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return result.PentestID, nil
}

// DownloadEvidence downloads and optionally unzips the evidence for a pentest
func (c *HTTPClient) DownloadEvidence(ctx context.Context, pentestID string, outputPath string, unzip bool) error {
	url := fmt.Sprintf("%s/api/v1/pentests/%s/download", c.BaseURL, pentestID)

	headers := map[string]string{
		"accept":     "application/zip",
		"api-key":    c.Credentials.APIKey,
		"api-secret": c.Credentials.APISecret,
	}

	zipFileName := fmt.Sprintf("pentest_%s_evidence.zip", pentestID)
	zipFilePath := filepath.Join(outputPath, zipFileName)

	err := helpers.DownloadFile(url, headers, zipFileName, outputPath, false)
	if err != nil {
		return fmt.Errorf("failed to download evidence: %w", err)
	}

	if unzip {
		extractDir := filepath.Join(outputPath, fmt.Sprintf("pentest_%s_evidence", pentestID))
		err = helpers.Unzip(zipFilePath, extractDir)
		if err != nil {
			return fmt.Errorf("failed to unzip evidence: %w", err)
		}
	}

	return nil
}
