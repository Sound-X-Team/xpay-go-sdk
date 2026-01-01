package xpay

import (
	"context"
	"net/http"
	"time"
)

// Config represents the configuration for the X-Pay client
type Config struct {
	APIKey      string
	MerchantID  string
	Environment Environment
	BaseURL     string
	Timeout     time.Duration
	HTTPClient  *http.Client
}

// Client is the main X-Pay client
type Client struct {
	config     *Config
	httpClient *HTTPClient
	merchantID string

	// API services
	Payments  *PaymentsService
	Customers *CustomersService
	Webhooks  *WebhooksService
}

// NewClient creates a new X-Pay client
func NewClient(config *Config) *Client {
	if config == nil {
		panic("config cannot be nil")
	}

	if config.APIKey == "" {
		panic("API key is required")
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.Environment == "" {
		// Auto-detect environment from API key
		if config.APIKey[:3] == "sk_" {
			if config.APIKey[3:11] == "sandbox_" {
				config.Environment = EnvironmentSandbox
			} else {
				config.Environment = EnvironmentLive
			}
		} else {
			config.Environment = EnvironmentSandbox
		}
	}

	// Extract merchant ID from config or use default
	merchantID := config.MerchantID
	if merchantID == "" {
		merchantID = "default" // In production, this would be extracted from API key or passed explicitly
	}

	httpClient := NewHTTPClient(config)

	client := &Client{
		config:     config,
		httpClient: httpClient,
		merchantID: merchantID,
	}

	// Initialize services
	client.Payments = NewPaymentsService(client)
	client.Customers = NewCustomersService(client)
	client.Webhooks = NewWebhooksService(client)

	return client
}

// Ping tests API connectivity and authentication
func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	var result APIResponse[PingData]
	err := c.httpClient.Get(ctx, "/v1/healthz", &result)
	if err != nil {
		return nil, err
	}

	return &PingResponse{
		Success:   result.Success,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// GetPaymentMethods retrieves available payment methods for the merchant
func (c *Client) GetPaymentMethods(ctx context.Context) (*PaymentMethods, error) {
	var result APIResponse[PaymentMethods]
	path := "/v1/api/merchants/" + c.merchantID + "/payment-methods"
	
	err := c.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// PingResponse represents the response from the ping endpoint
type PingResponse struct {
	Success   bool   `json:"success"`
	Timestamp string `json:"timestamp"`
}

// PingData represents the data from the ping endpoint
type PingData struct {
	Timestamp string `json:"timestamp"`
}