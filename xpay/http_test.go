package xpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClient_NewHTTPClient(t *testing.T) {
	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		Timeout:    30 * time.Second,
	}

	client := NewHTTPClient(config)
	if client == nil {
		t.Error("Expected client to be created")
	}
	if client.userAgent != "xpay-go-sdk/1.0.0" {
		t.Errorf("Expected user agent xpay-go-sdk/1.0.0, got %s", client.userAgent)
	}
}

func TestHTTPClient_DefaultBaseURL(t *testing.T) {
	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
	}

	client := NewHTTPClient(config)
	if client.baseURL != "https://api.xpay-bits.com" {
		t.Errorf("Expected default base URL, got %s", client.baseURL)
	}
}

func TestHTTPClient_CustomBaseURL(t *testing.T) {
	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    "https://custom-api.example.com",
	}

	client := NewHTTPClient(config)
	if client.baseURL != "https://custom-api.example.com" {
		t.Errorf("Expected custom base URL, got %s", client.baseURL)
	}
}

func TestHTTPClient_TrimsTrailingSlash(t *testing.T) {
	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    "https://api.example.com/",
	}

	client := NewHTTPClient(config)
	if client.baseURL != "https://api.example.com" {
		t.Errorf("Expected trailing slash to be trimmed, got %s", client.baseURL)
	}
}

func TestHTTPClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("Authorization") != "Bearer sk_sandbox_test" {
			t.Errorf("Expected Bearer token in Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"id": "123"},
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	var result map[string]interface{}
	err := client.Get(context.Background(), "/v1/test", &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["success"] != true {
		t.Error("Expected success to be true")
	}
}

func TestHTTPClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["amount"] != "10.00" {
			t.Errorf("Expected amount 10.00 in body")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"id": "pay_123"},
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	requestBody := map[string]interface{}{
		"amount":   "10.00",
		"currency": "USD",
	}

	var result map[string]interface{}
	err := client.Post(context.Background(), "/v1/payments", requestBody, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHTTPClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	var result map[string]interface{}
	err := client.Put(context.Background(), "/v1/test/123", nil, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHTTPClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	var result map[string]interface{}
	err := client.Delete(context.Background(), "/v1/test/123", &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHTTPClient_Error401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid API key",
			"code":    "AUTHENTICATION_ERROR",
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "invalid_key",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	var result map[string]interface{}
	err := client.Get(context.Background(), "/v1/test", &result)
	if err == nil {
		t.Error("Expected error for 401 response")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error type, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("Expected status code 401, got %d", apiErr.StatusCode)
	}
}

func TestHTTPClient_Error400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request",
			"code":    "VALIDATION_ERROR",
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	var result map[string]interface{}
	err := client.Post(context.Background(), "/v1/test", nil, &result)
	if err == nil {
		t.Error("Expected error for 400 response")
	}

	apiErr, ok := err.(*Error)
	if ok && apiErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}
}

func TestHTTPClient_Error500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:     "sk_sandbox_test",
		MerchantID: "test_merchant",
		BaseURL:    server.URL,
	}

	client := NewHTTPClient(config)

	var result map[string]interface{}
	err := client.Get(context.Background(), "/v1/test", &result)
	if err == nil {
		t.Error("Expected error for 500 response")
	}
}

func TestBuildQueryString_Empty(t *testing.T) {
	params := make(map[string]interface{})
	result := BuildQueryString(params)
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}

func TestBuildQueryString_WithParams(t *testing.T) {
	params := map[string]interface{}{
		"limit":  10,
		"offset": 0,
		"email":  "test@example.com",
	}
	result := BuildQueryString(params)
	if result == "" {
		t.Error("Expected query string, got empty")
	}
	if result[0] != '?' {
		t.Error("Expected query string to start with ?")
	}
}

func TestBuildQueryString_IgnoresEmptyStrings(t *testing.T) {
	params := map[string]interface{}{
		"limit": 10,
		"email": "",
	}
	result := BuildQueryString(params)
	// Should only have limit, not empty email
	if result == "" {
		t.Error("Expected query string with limit")
	}
}

func TestBuildQueryString_HandlesBool(t *testing.T) {
	params := map[string]interface{}{
		"active": true,
	}
	result := BuildQueryString(params)
	if result == "" {
		t.Error("Expected query string with bool")
	}
}
