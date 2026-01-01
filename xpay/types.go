package xpay

import (
	"time"
	
	"github.com/shopspring/decimal"
)

// Environment represents the X-Pay environment
type Environment string

const (
	EnvironmentSandbox Environment = "sandbox"
	EnvironmentLive    Environment = "live"
)

// PaymentMethod represents supported payment methods
type PaymentMethod string

const (
	PaymentMethodStripe      PaymentMethod = "stripe"
	PaymentMethodMomo        PaymentMethod = "momo"
	PaymentMethodMomoLiberia PaymentMethod = "momo_liberia"
	PaymentMethodMomoNigeria PaymentMethod = "momo_nigeria"
	PaymentMethodMomoUganda  PaymentMethod = "momo_uganda"
	PaymentMethodMomoRwanda  PaymentMethod = "momo_rwanda"
	PaymentMethodWallet      PaymentMethod = "wallet"
	PaymentMethodXPayWallet  PaymentMethod = "xpay_wallet"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing" 
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusSucceeded  PaymentStatus = "succeeded"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

// Currency represents supported currencies
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyGHS Currency = "GHS"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyNGN Currency = "NGN"
	CurrencyUGX Currency = "UGX"
	CurrencyRWF Currency = "RWF"
)

// CurrencyInfo represents currency metadata
type CurrencyInfo struct {
	Code             Currency `json:"code"`
	Name             string   `json:"name"`
	Symbol           string   `json:"symbol"`
	DecimalPlaces    int      `json:"decimal_places"`
	SmallestUnitName string   `json:"smallest_unit_name"`
}

// SupportedCurrencies maps currency codes to their information
var SupportedCurrencies = map[Currency]CurrencyInfo{
	CurrencyUSD: {
		Code:             CurrencyUSD,
		Name:             "US Dollar",
		Symbol:           "$",
		DecimalPlaces:    2,
		SmallestUnitName: "cents",
	},
	CurrencyGHS: {
		Code:             CurrencyGHS,
		Name:             "Ghanaian Cedi",
		Symbol:           "₵",
		DecimalPlaces:    2,
		SmallestUnitName: "pesewas",
	},
	CurrencyEUR: {
		Code:             CurrencyEUR,
		Name:             "Euro",
		Symbol:           "€",
		DecimalPlaces:    2,
		SmallestUnitName: "cents",
	},
	CurrencyGBP: {
		Code:             CurrencyGBP,
		Name:             "British Pound",
		Symbol:           "£",
		DecimalPlaces:    2,
		SmallestUnitName: "pence",
	},
}

// PaymentMethodData contains payment-method-specific data
type PaymentMethodData struct {
	// Stripe
	PaymentMethodTypes []string `json:"payment_method_types,omitempty"`
	
	// Mobile Money
	PhoneNumber string `json:"phone_number,omitempty"`
	
	// X-Pay Wallet
	WalletID string `json:"wallet_id,omitempty"`
	PIN      string `json:"pin,omitempty"`
}

// PaymentRequest represents a payment creation request
type PaymentRequest struct {
	Amount            decimal.Decimal            `json:"amount"`
	Currency          Currency                   `json:"currency,omitempty"`
	PaymentMethod     PaymentMethod              `json:"payment_method"`
	Description       string                     `json:"description,omitempty"`
	CustomerID        string                     `json:"customer_id,omitempty"`
	PaymentMethodData *PaymentMethodData         `json:"payment_method_data,omitempty"`
	Metadata          map[string]interface{}     `json:"metadata,omitempty"`
	SuccessURL        string                     `json:"success_url,omitempty"`
	CancelURL         string                     `json:"cancel_url,omitempty"`
	WebhookURL        string                     `json:"webhook_url,omitempty"`
}

// Payment represents a payment object
type Payment struct {
	ID              string                     `json:"id"`
	Status          PaymentStatus              `json:"status"`
	Amount          decimal.Decimal            `json:"amount"`
	Currency        Currency                   `json:"currency"`
	Description     string                     `json:"description,omitempty"`
	PaymentMethod   PaymentMethod              `json:"payment_method"`
	CustomerID      string                     `json:"customer_id,omitempty"`
	ClientSecret    string                     `json:"client_secret,omitempty"`
	ReferenceID     string                     `json:"reference_id,omitempty"`
	TransactionURL  string                     `json:"transaction_url,omitempty"`
	Instructions    string                     `json:"instructions,omitempty"`
	Metadata        map[string]interface{}     `json:"metadata,omitempty"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
}

// Customer represents a customer object
type Customer struct {
	ID        string                     `json:"id"`
	Name      string                     `json:"name"`
	Email     string                     `json:"email,omitempty"`
	Phone     string                     `json:"phone,omitempty"`
	Metadata  map[string]interface{}     `json:"metadata,omitempty"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
}

// CreateCustomerRequest represents a customer creation request
type CreateCustomerRequest struct {
	Name     string                     `json:"name"`
	Email    string                     `json:"email,omitempty"`
	Phone    string                     `json:"phone,omitempty"`
	Metadata map[string]interface{}     `json:"metadata,omitempty"`
}

// UpdateCustomerRequest represents a customer update request
type UpdateCustomerRequest struct {
	Name     string                     `json:"name,omitempty"`
	Email    string                     `json:"email,omitempty"`
	Phone    string                     `json:"phone,omitempty"`
	Metadata map[string]interface{}     `json:"metadata,omitempty"`
}

// ListCustomersRequest represents parameters for listing customers
type ListCustomersRequest struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
}

// WebhookEndpoint represents a webhook endpoint
type WebhookEndpoint struct {
	ID          string      `json:"id"`
	URL         string      `json:"url"`
	Events      []string    `json:"events"`
	Environment Environment `json:"environment"`
	IsActive    bool        `json:"is_active"`
	Secret      string      `json:"secret"`
	CreatedAt   time.Time   `json:"created_at"`
}

// CreateWebhookRequest represents a webhook creation request
type CreateWebhookRequest struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
}

// UpdateWebhookRequest represents a webhook update request
type UpdateWebhookRequest struct {
	URL         string   `json:"url,omitempty"`
	Events      []string `json:"events,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
	Description string   `json:"description,omitempty"`
}

// ListPaymentsRequest represents parameters for listing payments
type ListPaymentsRequest struct {
	Limit         int           `json:"limit,omitempty"`
	Offset        int           `json:"offset,omitempty"`
	Status        PaymentStatus `json:"status,omitempty"`
	CustomerID    string        `json:"customer_id,omitempty"`
	CreatedAfter  *time.Time    `json:"created_after,omitempty"`
	CreatedBefore *time.Time    `json:"created_before,omitempty"`
}

// PaymentMethodInfo represents information about a payment method
type PaymentMethodInfo struct {
	Type        string   `json:"type"`
	DisplayName string   `json:"display_name"`
	Currencies  []string `json:"currencies"`
	Description string   `json:"description"`
}

// PaymentMethods represents available payment methods
type PaymentMethods struct {
	AvailableMethods []PaymentMethodInfo `json:"available_methods"`
	Country          string              `json:"country"`
	DefaultCurrency  string              `json:"default_currency"`
}

// APIResponse represents a standard API response
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Error represents an X-Pay API error
type Error struct {
	Message    string                 `json:"message"`
	Code       string                 `json:"code"`
	StatusCode int                    `json:"status_code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}