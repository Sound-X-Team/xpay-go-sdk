package xpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomersService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api/merchants/test_merchant/customers" {
			t.Errorf("Expected correct path, got %s", r.URL.Path)
		}

		response := APIResponse[Customer]{
			Success: true,
			Data: Customer{
				ID:    "cust_123",
				Email: "test@example.com",
				Name:  "John Doe",
				Phone: "+1234567890",
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

	req := &CreateCustomerRequest{
		Email: "test@example.com",
		Name:  "John Doe",
		Phone: "+1234567890",
	}

	customer, err := client.Customers.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if customer.ID != "cust_123" {
		t.Errorf("Expected ID cust_123, got %s", customer.ID)
	}
	if customer.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", customer.Email)
	}
}

func TestCustomersService_Retrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api/merchants/test_merchant/customers/cust_123" {
			t.Errorf("Expected correct path, got %s", r.URL.Path)
		}

		response := APIResponse[Customer]{
			Success: true,
			Data: Customer{
				ID:    "cust_123",
				Email: "test@example.com",
				Name:  "John Doe",
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

	customer, err := client.Customers.Retrieve(context.Background(), "cust_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if customer.ID != "cust_123" {
		t.Errorf("Expected ID cust_123, got %s", customer.ID)
	}
}

func TestCustomersService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}

		response := APIResponse[Customer]{
			Success: true,
			Data: Customer{
				ID:    "cust_123",
				Email: "test@example.com",
				Name:  "John Updated",
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

	req := &UpdateCustomerRequest{
		Name: "John Updated",
	}

	customer, err := client.Customers.Update(context.Background(), "cust_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if customer.Name != "John Updated" {
		t.Errorf("Expected name 'John Updated', got %s", customer.Name)
	}
}

func TestCustomersService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api/merchants/test_merchant/customers/cust_123" {
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

	err := client.Customers.Delete(context.Background(), "cust_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCustomersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		response := APIResponse[struct {
			Customers []Customer `json:"customers"`
			Total     int        `json:"total"`
		}]{
			Success: true,
			Data: struct {
				Customers []Customer `json:"customers"`
				Total     int        `json:"total"`
			}{
				Customers: []Customer{
					{ID: "cust_1", Email: "user1@example.com", Name: "User 1"},
					{ID: "cust_2", Email: "user2@example.com", Name: "User 2"},
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

	customers, err := client.Customers.List(context.Background(), &ListCustomersRequest{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(customers.Items) != 2 {
		t.Errorf("Expected 2 customers, got %d", len(customers.Items))
	}
	if customers.Total != 2 {
		t.Errorf("Expected total 2, got %d", customers.Total)
	}
}

func TestCustomersService_ListWithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check query parameters
		email := r.URL.Query().Get("email")
		if email != "test@example.com" {
			t.Errorf("Expected email filter, got %s", email)
		}

		response := APIResponse[struct {
			Customers []Customer `json:"customers"`
			Total     int        `json:"total"`
		}]{
			Success: true,
			Data: struct {
				Customers []Customer `json:"customers"`
				Total     int        `json:"total"`
			}{
				Customers: []Customer{
					{ID: "cust_1", Email: "test@example.com", Name: "Test User"},
				},
				Total: 1,
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

	customers, err := client.Customers.List(context.Background(), &ListCustomersRequest{
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(customers.Items) != 1 {
		t.Errorf("Expected 1 customer, got %d", len(customers.Items))
	}
}
