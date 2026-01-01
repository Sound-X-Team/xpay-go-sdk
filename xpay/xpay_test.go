package xpay

import (
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
)

func TestCurrencyUtilities(t *testing.T) {
	// Test ToSmallestUnit
	amount := decimal.NewFromFloat(25.99)
	cents := ToSmallestUnit(amount, CurrencyUSD)
	expected := int64(2599)
	
	if cents != expected {
		t.Errorf("ToSmallestUnit failed: expected %d, got %d", expected, cents)
	}
	
	// Test FromSmallestUnit
	dollarsBack := FromSmallestUnit(cents, CurrencyUSD)
	if !dollarsBack.Equal(amount) {
		t.Errorf("FromSmallestUnit failed: expected %s, got %s", amount.String(), dollarsBack.String())
	}
	
	// Test FormatAmount
	formatted := FormatAmount(cents, CurrencyUSD, true)
	expectedFormat := "$25.99"
	
	if formatted != expectedFormat {
		t.Errorf("FormatAmount failed: expected %s, got %s", expectedFormat, formatted)
	}
}

func TestPaymentMethodCurrencies(t *testing.T) {
	// Test Stripe supported currencies
	stripeConfig := PaymentMethodCurrencies[PaymentMethodStripe]
	expectedCurrencies := []Currency{CurrencyUSD, CurrencyEUR, CurrencyGBP, CurrencyGHS}
	
	if len(stripeConfig.SupportedCurrencies) != len(expectedCurrencies) {
		t.Errorf("Stripe currency count mismatch: expected %d, got %d", 
			len(expectedCurrencies), len(stripeConfig.SupportedCurrencies))
	}
	
	// Test default currency
	if stripeConfig.DefaultCurrency != CurrencyUSD {
		t.Errorf("Stripe default currency mismatch: expected %s, got %s", 
			CurrencyUSD, stripeConfig.DefaultCurrency)
	}
}

func TestErrorStructure(t *testing.T) {
	err := &Error{
		Message:    "Test error",
		Code:       "TEST_ERROR",
		StatusCode: 400,
		Details:    map[string]interface{}{"field": "value"},
	}
	
	if err.Error() != "Test error" {
		t.Errorf("Error.Error() failed: expected 'Test error', got '%s'", err.Error())
	}
}

func TestConfigDefaults(t *testing.T) {
	config := &Config{
		APIKey:     "sk_sandbox_test_key",
		MerchantID: "test_merchant",
	}
	
	client := NewClient(config)
	
	if client.config.Environment != EnvironmentSandbox {
		t.Errorf("Environment auto-detection failed: expected %s, got %s", 
			EnvironmentSandbox, client.config.Environment)
	}
	
	if client.config.Timeout != 30*time.Second {
		t.Errorf("Default timeout failed: expected %v, got %v", 
			30*time.Second, client.config.Timeout)
	}
}

func TestPaymentRequestValidation(t *testing.T) {
	// Test that PaymentRequest struct is properly defined
	req := &PaymentRequest{
		Amount:        decimal.NewFromFloat(10.00),
		Currency:      CurrencyUSD,
		PaymentMethod: PaymentMethodStripe,
		Description:   "Test payment",
		PaymentMethodData: &PaymentMethodData{
			PaymentMethodTypes: []string{"card"},
		},
		Metadata: map[string]interface{}{
			"test": true,
		},
	}
	
	if req.Amount.IsZero() {
		t.Error("PaymentRequest amount should not be zero")
	}
	
	if req.Currency != CurrencyUSD {
		t.Errorf("Currency mismatch: expected %s, got %s", CurrencyUSD, req.Currency)
	}
}

func TestAPIResponseStructure(t *testing.T) {
	// Test generic API response
	response := APIResponse[Payment]{
		Success: true,
		Data: Payment{
			ID:     "pay_123",
			Status: PaymentStatusSucceeded,
			Amount: decimal.NewFromFloat(25.99),
		},
		Message: "Success",
	}
	
	if !response.Success {
		t.Error("APIResponse should be successful")
	}
	
	if response.Data.ID != "pay_123" {
		t.Errorf("Payment ID mismatch: expected 'pay_123', got '%s'", response.Data.ID)
	}
}

func TestWebhookSignatureVerification(t *testing.T) {
	// Test the signature verification logic exists and is callable
	payload := []byte(`{"event": "payment.succeeded"}`)
	signature := "test_signature"
	secret := "test_secret"
	
	// This should not panic and should return a boolean
	result := VerifyWebhookSignature(payload, signature, secret)
	
	// Since we're using a test signature, it should be false
	if result {
		t.Error("Signature verification should fail with test data")
	}
}