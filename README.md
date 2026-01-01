# X-Pay Go SDK

[![Go Version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://golang.org/dl/)
[![GoDoc](https://godoc.org/github.com/Sound-X-Team/xpay-go-sdk?status.svg)](https://godoc.org/github.com/Sound-X-Team/xpay-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Go SDK for X-Pay payment processing platform. Accept payments from multiple providers including Stripe, Mobile Money, and X-Pay Wallets with a unified API.

## Features

- 🚀 **Multiple Payment Methods**: Stripe, Mobile Money (MoMo), X-Pay Wallets
- 🔒 **Type Safety**: Full Go struct definitions with validation tags
- 💰 **Decimal Precision**: Proper currency handling with shopspring/decimal
- 🌍 **Environment Detection**: Auto-detect sandbox/live from API keys
- 👥 **Customer Management**: Complete customer lifecycle management
- 🔗 **Webhook Management**: Easy webhook setup and verification
- ⚡ **Context Support**: Full context.Context support for cancellation and timeouts
- 🧪 **Testing**: Comprehensive test suite with mock clients

## Installation

```bash
go get github.com/Sound-X-Team/xpay-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Sound-X-Team/xpay-go-sdk/xpay"
    "github.com/shopspring/decimal"
)

func main() {
    // Initialize client
    client := xpay.NewClient(&xpay.Config{
        APIKey:      "sk_sandbox_your_secret_key_here",
        MerchantID:  "your_merchant_id_here",
        Environment: xpay.EnvironmentSandbox,
        BaseURL:     "http://localhost:8000", // For local development
    })

    // Create a payment
    payment, err := client.Payments.Create(context.Background(), &xpay.PaymentRequest{
        Amount:        decimal.NewFromFloat(29.99),
        Currency:      "USD",
        PaymentMethod: xpay.PaymentMethodStripe,
        Description:   "Go SDK test payment",
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Payment created: %s\n", payment.ID)
    fmt.Printf("Status: %s\n", payment.Status)
}
```

## Working Credentials (for testing)

```go
client := xpay.NewClient(&xpay.Config{
    APIKey:      "sk_sandbox_7c845adf-f658-4f29-9857-7e8a8708",
    MerchantID:  "548d8033-fbe9-411b-991f-f159cdee7745",
    Environment: xpay.EnvironmentSandbox,
    BaseURL:     "http://localhost:8000",
})
```

## Payment Methods

All payment methods from other SDKs are supported:

- `stripe` - Credit/debit cards
- `momo_liberia` - MTN Mobile Money Liberia  
- `xpay_wallet` - X-Pay wallet transfers
- More coming soon...

## Usage Examples

### Stripe Payment

```go
payment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
    Amount:        decimal.NewFromFloat(25.99),
    Currency:      "USD",
    PaymentMethod: xpay.PaymentMethodStripe,
    Description:   "Stripe test payment",
    PaymentMethodData: &xpay.PaymentMethodData{
        PaymentMethodTypes: []string{"card"},
    },
})
```

### Mobile Money Payment

```go
payment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
    Amount:        decimal.NewFromFloat(50.00),
    Currency:      "USD",
    PaymentMethod: xpay.PaymentMethodMomoLiberia,
    Description:   "MoMo payment",
    PaymentMethodData: &xpay.PaymentMethodData{
        PhoneNumber: "+231700000000",
    },
})
```

### X-Pay Wallet Payment

```go
payment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
    Amount:        decimal.NewFromFloat(15.50),
    Currency:      "USD",
    PaymentMethod: xpay.PaymentMethodXPayWallet,
    Description:   "Wallet transfer",
    PaymentMethodData: &xpay.PaymentMethodData{
        WalletID: "wallet_123",
        PIN:      "1234",
    },
})
```

### Customer Management

```go
// Create a customer
customer, err := client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
    Name:  "John Doe",
    Email: "john@example.com",
    Phone: "+231700000000",
    Metadata: map[string]interface{}{
        "source": "website",
    },
})

// Retrieve a customer
customer, err := client.Customers.Retrieve(ctx, "cust_123")

// List customers
customers, err := client.Customers.List(ctx, &xpay.ListCustomersRequest{
    Limit: 10,
    Email: "john@example.com",
})
```

### Webhook Management

```go
// Create a webhook endpoint
webhook, err := client.Webhooks.Create(ctx, &xpay.CreateWebhookRequest{
    URL:    "https://your-site.com/webhooks/xpay",
    Events: []string{"payment.succeeded", "payment.failed"},
})

// Verify webhook signature
isValid := xpay.VerifyWebhookSignature(
    payload,
    signature,
    "whsec_your_webhook_secret",
)
```

## Error Handling

```go
payment, err := client.Payments.Create(ctx, paymentRequest)
if err != nil {
    var xpayErr *xpay.Error
    if errors.As(err, &xpayErr) {
        fmt.Printf("X-Pay error [%s]: %s\n", xpayErr.Code, xpayErr.Message)
        if xpayErr.Details != nil {
            fmt.Printf("Details: %+v\n", xpayErr.Details)
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Currency Utilities

```go
import "github.com/shopspring/decimal"

// Convert to smallest unit (cents)
amount := decimal.NewFromFloat(25.99)
cents := xpay.ToSmallestUnit(amount, "USD") // 2599

// Convert from smallest unit
dollars := xpay.FromSmallestUnit(cents, "USD") // 25.99

// Format for display
formatted := xpay.FormatAmount(cents, "USD", true) // "$25.99"
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires valid API keys)
go test -tags=integration ./...
```

## Configuration

### Environment Variables

```bash
export XPAY_API_KEY="sk_sandbox_your_key_here"
export XPAY_MERCHANT_ID="your_merchant_id_here"
export XPAY_ENVIRONMENT="sandbox"
export XPAY_BASE_URL="http://localhost:8000"
export XPAY_TIMEOUT="30s"
```

### Configuration Struct

```go
config := &xpay.Config{
    APIKey:      "sk_sandbox_your_key_here",
    MerchantID:  "your_merchant_id_here",
    Environment: xpay.EnvironmentSandbox,
    BaseURL:     "http://localhost:8000",
    Timeout:     30 * time.Second,
    HTTPClient:  &http.Client{}, // Optional custom HTTP client
}

client := xpay.NewClient(config)
```

## Requirements

- Go 1.19+
- github.com/shopspring/decimal
- github.com/google/uuid

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)  
5. Open Pull Request

## License

MIT License - see LICENSE file for details.