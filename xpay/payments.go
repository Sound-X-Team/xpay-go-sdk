package xpay

import (
	"context"
	"fmt"
	
	"github.com/shopspring/decimal"
)

// PaymentsService handles payment-related operations
type PaymentsService struct {
	client *Client
}

// NewPaymentsService creates a new payments service
func NewPaymentsService(client *Client) *PaymentsService {
	return &PaymentsService{client: client}
}

// PaymentMethodCurrencies maps payment methods to their supported currencies
var PaymentMethodCurrencies = map[PaymentMethod]PaymentMethodCurrency{
	PaymentMethodStripe: {
		PaymentMethod:        PaymentMethodStripe,
		SupportedCurrencies:  []Currency{CurrencyUSD, CurrencyEUR, CurrencyGBP, CurrencyGHS},
		DefaultCurrency:      CurrencyUSD,
		Regions:              []string{"US", "EU", "GB", "GH"},
	},
	PaymentMethodMomo: {
		PaymentMethod:        PaymentMethodMomo,
		SupportedCurrencies:  []Currency{CurrencyGHS},
		DefaultCurrency:      CurrencyGHS,
		Regions:              []string{"GH"},
	},
	PaymentMethodMomoLiberia: {
		PaymentMethod:        PaymentMethodMomoLiberia,
		SupportedCurrencies:  []Currency{CurrencyUSD},
		DefaultCurrency:      CurrencyUSD,
		Regions:              []string{"LR"},
	},
	PaymentMethodXPayWallet: {
		PaymentMethod:        PaymentMethodXPayWallet,
		SupportedCurrencies:  []Currency{CurrencyUSD, CurrencyGHS, CurrencyEUR},
		DefaultCurrency:      CurrencyUSD,
		Regions:              []string{"US", "GH", "EU"},
	},
}

// PaymentMethodCurrency represents currency information for a payment method
type PaymentMethodCurrency struct {
	PaymentMethod       PaymentMethod `json:"payment_method"`
	SupportedCurrencies []Currency    `json:"supported_currencies"`
	DefaultCurrency     Currency      `json:"default_currency"`
	Regions             []string      `json:"regions"`
}

// getDefaultCurrency returns the default currency for a payment method
func (s *PaymentsService) getDefaultCurrency(paymentMethod PaymentMethod) Currency {
	if methodConfig, exists := PaymentMethodCurrencies[paymentMethod]; exists {
		return methodConfig.DefaultCurrency
	}
	return CurrencyUSD
}

// validateCurrency validates that a currency is supported for a payment method
func (s *PaymentsService) validateCurrency(paymentMethod PaymentMethod, currency Currency) error {
	methodConfig, exists := PaymentMethodCurrencies[paymentMethod]
	if !exists {
		return &Error{
			Message: fmt.Sprintf("Unsupported payment method: %s", paymentMethod),
			Code:    "INVALID_PAYMENT_METHOD",
		}
	}

	for _, supportedCurrency := range methodConfig.SupportedCurrencies {
		if supportedCurrency == currency {
			return nil
		}
	}

	return &Error{
		Message: fmt.Sprintf("Currency %s is not supported for payment method %s. Supported currencies: %v", 
			currency, paymentMethod, methodConfig.SupportedCurrencies),
		Code: "INVALID_CURRENCY",
	}
}

// Create creates a new payment
func (s *PaymentsService) Create(ctx context.Context, req *PaymentRequest) (*Payment, error) {
	// Auto-assign currency if not provided
	processedReq := *req
	if processedReq.Currency == "" {
		processedReq.Currency = s.getDefaultCurrency(req.PaymentMethod)
	}

	// Validate currency for the payment method
	if err := s.validateCurrency(req.PaymentMethod, processedReq.Currency); err != nil {
		return nil, err
	}

	var result APIResponse[Payment]
	path := fmt.Sprintf("/v1/api/merchants/%s/payments", s.client.merchantID)
	
	err := s.client.httpClient.Post(ctx, path, &processedReq, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Retrieve retrieves a payment by ID
func (s *PaymentsService) Retrieve(ctx context.Context, paymentID string) (*Payment, error) {
	var result APIResponse[Payment]
	path := fmt.Sprintf("/v1/api/merchants/%s/payments/%s", s.client.merchantID, paymentID)
	
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// List lists payments with optional filters
func (s *PaymentsService) List(ctx context.Context, req *ListPaymentsRequest) (*ListResponse[Payment], error) {
	params := make(map[string]interface{})
	
	if req != nil {
		if req.Limit > 0 {
			params["limit"] = req.Limit
		}
		if req.Offset > 0 {
			params["offset"] = req.Offset
		}
		if req.Status != "" {
			params["status"] = string(req.Status)
		}
		if req.CustomerID != "" {
			params["customer_id"] = req.CustomerID
		}
		if req.CreatedAfter != nil {
			params["created_after"] = req.CreatedAfter
		}
		if req.CreatedBefore != nil {
			params["created_before"] = req.CreatedBefore
		}
	}

	queryString := BuildQueryString(params)
	path := fmt.Sprintf("/v1/api/merchants/%s/payments%s", s.client.merchantID, queryString)

	var result APIResponse[struct {
		Payments []Payment `json:"payments"`
		Total    int       `json:"total"`
	}]
	
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &ListResponse[Payment]{
		Items:  result.Data.Payments,
		Total:  result.Data.Total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// Cancel cancels a payment (if supported by payment method)
func (s *PaymentsService) Cancel(ctx context.Context, paymentID string) (*Payment, error) {
	var result APIResponse[Payment]
	path := fmt.Sprintf("/v1/api/merchants/%s/payments/%s/cancel", s.client.merchantID, paymentID)
	
	err := s.client.httpClient.Post(ctx, path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Confirm confirms a payment (for payment methods that require confirmation)
func (s *PaymentsService) Confirm(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*Payment, error) {
	var result APIResponse[Payment]
	path := fmt.Sprintf("/v1/payments/%s/confirm", paymentID)
	
	err := s.client.httpClient.Post(ctx, path, confirmationData, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetPaymentMethods gets available payment methods for this merchant
func (s *PaymentsService) GetPaymentMethods(ctx context.Context, country string) (*PaymentMethods, error) {
	params := make(map[string]interface{})
	if country != "" {
		params["country"] = country
	}

	queryString := BuildQueryString(params)
	path := fmt.Sprintf("/v1/api/merchants/%s/payment-methods%s", s.client.merchantID, queryString)

	var result APIResponse[PaymentMethods]
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetSupportedCurrencies returns supported currencies for a payment method
func (s *PaymentsService) GetSupportedCurrencies(paymentMethod PaymentMethod) []Currency {
	if methodConfig, exists := PaymentMethodCurrencies[paymentMethod]; exists {
		return methodConfig.SupportedCurrencies
	}
	return []Currency{}
}

// ToSmallestUnit converts an amount to the smallest currency unit (e.g., dollars to cents)
func ToSmallestUnit(amount decimal.Decimal, currency Currency) int64 {
	currencyInfo, exists := SupportedCurrencies[currency]
	if !exists {
		panic(fmt.Sprintf("Unsupported currency: %s", currency))
	}
	
	multiplier := decimal.NewFromInt(int64(pow10(currencyInfo.DecimalPlaces)))
	result := amount.Mul(multiplier)
	return result.IntPart()
}

// FromSmallestUnit converts an amount from the smallest currency unit (e.g., cents to dollars)
func FromSmallestUnit(amount int64, currency Currency) decimal.Decimal {
	currencyInfo, exists := SupportedCurrencies[currency]
	if !exists {
		panic(fmt.Sprintf("Unsupported currency: %s", currency))
	}
	
	divisor := decimal.NewFromInt(int64(pow10(currencyInfo.DecimalPlaces)))
	return decimal.NewFromInt(amount).Div(divisor)
}

// FormatAmount formats an amount for display with currency symbol
func FormatAmount(amount int64, currency Currency, fromSmallestUnit bool) string {
	currencyInfo, exists := SupportedCurrencies[currency]
	if !exists {
		panic(fmt.Sprintf("Unsupported currency: %s", currency))
	}
	
	var displayAmount decimal.Decimal
	if fromSmallestUnit {
		displayAmount = FromSmallestUnit(amount, currency)
	} else {
		displayAmount = decimal.NewFromInt(amount)
	}
	
	formatted := displayAmount.StringFixed(int32(currencyInfo.DecimalPlaces))
	return currencyInfo.Symbol + formatted
}

// pow10 calculates 10^n
func pow10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}