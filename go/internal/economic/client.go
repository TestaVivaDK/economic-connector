package economic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/TestaVivaDK/e-conomic-connector/internal/logger"
)

const apiBase = "https://restapi.e-conomic.com"

// Client is a thin wrapper around the e-conomic REST API.
type Client struct {
	baseURL             string
	appSecretToken      string
	agreementGrantToken string
	http                *http.Client
}

// NewClient creates an e-conomic API client.
func NewClient(appSecretToken, agreementGrantToken string) *Client {
	return &Client{
		baseURL:             apiBase,
		appSecretToken:      appSecretToken,
		agreementGrantToken: agreementGrantToken,
		http:                &http.Client{Timeout: 30 * time.Second},
	}
}

// buildURL constructs a full URL from a path and optional query parameters.
func (c *Client) buildURL(path string, queryParams map[string]string) string {
	u := c.baseURL + path
	if len(queryParams) == 0 {
		return u
	}
	var parts []string
	for k, v := range queryParams {
		parts = append(parts, k+"="+url.QueryEscape(v))
	}
	sep := "?"
	if strings.Contains(u, "?") {
		sep = "&"
	}
	return u + sep + strings.Join(parts, "&")
}

// Get performs a GET request.
func (c *Client) Get(path string, queryParams map[string]string) (json.RawMessage, error) {
	reqURL := c.buildURL(path, queryParams)
	return c.doRequest(http.MethodGet, reqURL, nil)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, body any) (json.RawMessage, error) {
	reqURL := c.baseURL + path
	return c.doRequest(http.MethodPost, reqURL, body)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, body any) (json.RawMessage, error) {
	reqURL := c.baseURL + path
	return c.doRequest(http.MethodPut, reqURL, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) (json.RawMessage, error) {
	reqURL := c.baseURL + path
	return c.doRequest(http.MethodDelete, reqURL, nil)
}

// TestConnection calls /self to verify the tokens are valid.
func (c *Client) TestConnection() (json.RawMessage, error) {
	return c.Get("/self", nil)
}

func (c *Client) doRequest(method, reqURL string, body any) (json.RawMessage, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-AppSecretToken", c.appSecretToken)
	req.Header.Set("X-AgreementGrantToken", c.agreementGrantToken)
	req.Header.Set("Content-Type", "application/json")

	if logger.Log != nil {
		logger.Log.Info(fmt.Sprintf("[ECONOMIC] %s %s", method, reqURL))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		message := fmt.Sprintf("e-conomic API error %d", resp.StatusCode)
		var parsed struct {
			Message string `json:"message"`
			Errors  any    `json:"errors"`
		}
		if json.Unmarshal(respBody, &parsed) == nil && parsed.Message != "" {
			message = parsed.Message
		}
		return nil, fmt.Errorf("%s", message)
	}

	if len(respBody) == 0 {
		return json.RawMessage(`{"success":true}`), nil
	}

	return json.RawMessage(respBody), nil
}
