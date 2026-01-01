package xpay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// WebhooksService handles webhook-related operations
type WebhooksService struct {
	client *Client
}

// NewWebhooksService creates a new webhooks service
func NewWebhooksService(client *Client) *WebhooksService {
	return &WebhooksService{client: client}
}

// Create creates a new webhook endpoint
func (s *WebhooksService) Create(ctx context.Context, req *CreateWebhookRequest) (*WebhookEndpoint, error) {
	var result APIResponse[WebhookEndpoint]
	path := fmt.Sprintf("/v1/api/merchants/%s/webhooks", s.client.merchantID)
	
	err := s.client.httpClient.Post(ctx, path, req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Retrieve retrieves a webhook endpoint by ID
func (s *WebhooksService) Retrieve(ctx context.Context, webhookID string) (*WebhookEndpoint, error) {
	var result APIResponse[WebhookEndpoint]
	path := fmt.Sprintf("/v1/api/merchants/%s/webhooks/%s", s.client.merchantID, webhookID)
	
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Update updates a webhook endpoint
func (s *WebhooksService) Update(ctx context.Context, webhookID string, req *UpdateWebhookRequest) (*WebhookEndpoint, error) {
	var result APIResponse[WebhookEndpoint]
	path := fmt.Sprintf("/v1/api/merchants/%s/webhooks/%s", s.client.merchantID, webhookID)
	
	err := s.client.httpClient.Put(ctx, path, req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// List lists webhook endpoints
func (s *WebhooksService) List(ctx context.Context) (*ListResponse[WebhookEndpoint], error) {
	path := fmt.Sprintf("/v1/api/merchants/%s/webhooks", s.client.merchantID)

	var result APIResponse[struct {
		Webhooks []WebhookEndpoint `json:"webhooks"`
		Total    int               `json:"total"`
	}]
	
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &ListResponse[WebhookEndpoint]{
		Items: result.Data.Webhooks,
		Total: result.Data.Total,
	}, nil
}

// Delete deletes a webhook endpoint
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	path := fmt.Sprintf("/v1/api/merchants/%s/webhooks/%s", s.client.merchantID, webhookID)
	
	var result APIResponse[interface{}]
	err := s.client.httpClient.Delete(ctx, path, &result)
	if err != nil {
		return err
	}

	return nil
}

// Test tests a webhook endpoint
func (s *WebhooksService) Test(ctx context.Context, webhookID string) (*TestWebhookResponse, error) {
	var result APIResponse[TestWebhookResponse]
	path := fmt.Sprintf("/v1/api/merchants/%s/webhooks/%s/test", s.client.merchantID, webhookID)
	
	err := s.client.httpClient.Post(ctx, path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// TestWebhookResponse represents the response from testing a webhook
type TestWebhookResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// VerifyWebhookSignature verifies the signature of a webhook payload
func VerifyWebhookSignature(payload []byte, signature string, secret string) bool {
	// Remove "sha256=" prefix if present
	signature = strings.TrimPrefix(signature, "sha256=")
	
	// Create HMAC-SHA256 hash
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	
	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}