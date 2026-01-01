package xpay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhooksService_Create(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api/merchants/test_merchant/webhooks" {
			t.Errorf("Expected correct path, got %s", r.URL.Path)
		}

		response := APIResponse[WebhookEndpoint]{
			Success: true,
			Data: WebhookEndpoint{
				ID:          "webhook_123",
				URL:         "https://example.com/webhooks",
				Events:      []string{"payment.succeeded"},
				Environment: "sandbox",
				IsActive:    true,
				Secret:      "whsec_test",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	})

	req := &CreateWebhookRequest{
		URL:    "https://example.com/webhooks",
		Events: []string{"payment.succeeded"},
	}

	webhook, err := client.Webhooks.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if webhook.ID != "webhook_123" {
		t.Errorf("Expected ID webhook_123, got %s", webhook.ID)
	}
}

func TestWebhooksService_Retrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api/merchants/test_merchant/webhooks/webhook_123" {
			t.Errorf("Expected correct path, got %s", r.URL.Path)
		}

		response := APIResponse[WebhookEndpoint]{
			Success: true,
			Data: WebhookEndpoint{
				ID:          "webhook_123",
				URL:         "https://example.com/webhooks",
				Events:      []string{"payment.succeeded"},
				Environment: "sandbox",
				IsActive:    true,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	})

	webhook, err := client.Webhooks.Retrieve(context.Background(), "webhook_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if webhook.ID != "webhook_123" {
		t.Errorf("Expected ID webhook_123, got %s", webhook.ID)
	}
}

func TestWebhooksService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		response := APIResponse[struct {
			Webhooks []WebhookEndpoint `json:"webhooks"`
			Total    int               `json:"total"`
		}]{
			Success: true,
			Data: struct {
				Webhooks []WebhookEndpoint `json:"webhooks"`
				Total    int               `json:"total"`
			}{
				Webhooks: []WebhookEndpoint{
					{ID: "webhook_1", URL: "https://example.com/hook1"},
					{ID: "webhook_2", URL: "https://example.com/hook2"},
				},
				Total: 2,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	})

	webhooks, err := client.Webhooks.List(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(webhooks.Items) != 2 {
		t.Errorf("Expected 2 webhooks, got %d", len(webhooks.Items))
	}
	if webhooks.Total != 2 {
		t.Errorf("Expected total 2, got %d", webhooks.Total)
	}
}

func TestWebhooksService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api/merchants/test_merchant/webhooks/webhook_123" {
			t.Errorf("Expected correct path, got %s", r.URL.Path)
		}

		response := APIResponse[interface{}]{Success: true}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	})

	err := client.Webhooks.Delete(context.Background(), "webhook_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestWebhooksService_Test(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		response := APIResponse[TestWebhookResponse]{
			Success: true,
			Data: TestWebhookResponse{
				Success:   true,
				Message:   "Webhook test successful",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	})

	result, err := client.Webhooks.Test(context.Background(), "webhook_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}
}

// Signature verification tests
func TestVerifyWebhookSignature_Valid(t *testing.T) {
	payload := []byte(`{"event": "payment.succeeded", "data": {"id": "pay_123"}}`)
	secret := "whsec_test_secret"

	// Generate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	result := VerifyWebhookSignature(payload, signature, secret)
	if !result {
		t.Error("Expected signature to be valid")
	}
}

func TestVerifyWebhookSignature_WithPrefix(t *testing.T) {
	payload := []byte(`{"event": "payment.succeeded"}`)
	secret := "whsec_test_secret"

	// Generate valid signature with prefix
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	result := VerifyWebhookSignature(payload, signature, secret)
	if !result {
		t.Error("Expected signature with prefix to be valid")
	}
}

func TestVerifyWebhookSignature_Invalid(t *testing.T) {
	payload := []byte(`{"event": "payment.succeeded"}`)
	secret := "whsec_test_secret"
	invalidSignature := "invalid_signature_here"

	result := VerifyWebhookSignature(payload, invalidSignature, secret)
	if result {
		t.Error("Expected signature to be invalid")
	}
}

func TestVerifyWebhookSignature_TamperedPayload(t *testing.T) {
	originalPayload := []byte(`{"event": "payment.succeeded"}`)
	tamperedPayload := []byte(`{"event": "payment.failed"}`)
	secret := "whsec_test_secret"

	// Generate signature for original payload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(originalPayload)
	signature := hex.EncodeToString(mac.Sum(nil))

	// Verify with tampered payload should fail
	result := VerifyWebhookSignature(tamperedPayload, signature, secret)
	if result {
		t.Error("Expected tampered payload verification to fail")
	}
}
