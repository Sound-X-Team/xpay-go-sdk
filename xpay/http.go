package xpay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// HTTPClient handles HTTP requests to the X-Pay API
type HTTPClient struct {
	config     *Config
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(config *Config) *HTTPClient {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.xpay-bits.com"
	}

	return &HTTPClient{
		config:     config,
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		userAgent:  "xpay-go-sdk/1.0.0",
	}
}

// Get performs a GET request
func (c *HTTPClient) Get(ctx context.Context, path string, result interface{}) error {
	return c.makeRequest(ctx, "GET", path, nil, result)
}

// Post performs a POST request
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.makeRequest(ctx, "POST", path, body, result)
}

// Put performs a PUT request
func (c *HTTPClient) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.makeRequest(ctx, "PUT", path, body, result)
}

// Delete performs a DELETE request
func (c *HTTPClient) Delete(ctx context.Context, path string, result interface{}) error {
	return c.makeRequest(ctx, "DELETE", path, nil, result)
}

// makeRequest performs the actual HTTP request
func (c *HTTPClient) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add merchant ID if available
	if c.config.MerchantID != "" {
		req.Header.Set("X-Merchant-ID", c.config.MerchantID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var apiErr Error
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			// If we can't parse the error, create a generic one
			return &Error{
				Message:    string(respBody),
				Code:       "UNKNOWN_ERROR",
				StatusCode: resp.StatusCode,
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	// Parse successful response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// BuildQueryString builds a query string from parameters
func BuildQueryString(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for key, value := range params {
		if value == nil {
			continue
		}

		switch v := value.(type) {
		case string:
			if v != "" {
				values.Add(key, v)
			}
		case int:
			if v != 0 {
				values.Add(key, strconv.Itoa(v))
			}
		case bool:
			values.Add(key, strconv.FormatBool(v))
		case time.Time:
			values.Add(key, v.Format(time.RFC3339))
		case *time.Time:
			if v != nil {
				values.Add(key, v.Format(time.RFC3339))
			}
		default:
			// Convert to string for other types
			values.Add(key, fmt.Sprintf("%v", v))
		}
	}

	if len(values) == 0 {
		return ""
	}

	return "?" + values.Encode()
}