# Getting Started with X-Pay Go SDK

This guide will help you get started with the X-Pay Go SDK.

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
    // Initialize the client
    client := xpay.NewClient(&xpay.Config{
        APIKey:      "sk_sandbox_your_secret_key_here",
        MerchantID:  "your_merchant_id_here",
        Environment: xpay.EnvironmentSandbox,
        BaseURL:     "http://localhost:8000", // For local development
    })

    ctx := context.Background()

    // Create a payment
    payment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
        Amount:        decimal.NewFromFloat(29.99),
        Currency:      xpay.CurrencyUSD,
        PaymentMethod: xpay.PaymentMethodStripe,
        Description:   "Test payment",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment created: %s\n", payment.ID)
}
```

## Configuration

The SDK can be configured using environment variables or the Config struct:

### Environment Variables

```bash
export XPAY_API_KEY="sk_sandbox_your_key_here"
export XPAY_MERCHANT_ID="your_merchant_id_here"
export XPAY_ENVIRONMENT="sandbox"
export XPAY_BASE_URL="http://localhost:8000"
```

### Config Struct

```go
config := &xpay.Config{
    APIKey:      "sk_sandbox_your_key_here",
    MerchantID:  "your_merchant_id_here",
    Environment: xpay.EnvironmentSandbox,
    BaseURL:     "http://localhost:8000",
    Timeout:     30 * time.Second,
}

client := xpay.NewClient(config)
```

## Working with Payments

### Create a Payment

```go
payment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
    Amount:        decimal.NewFromFloat(25.99),
    Currency:      xpay.CurrencyUSD,
    PaymentMethod: xpay.PaymentMethodStripe,
    Description:   "Test payment",
    PaymentMethodData: &xpay.PaymentMethodData{
        PaymentMethodTypes: []string{"card"},
    },
})
```

### Retrieve a Payment

```go
payment, err := client.Payments.Retrieve(ctx, "pay_123456789")
```

### List Payments

```go
payments, err := client.Payments.List(ctx, &xpay.ListPaymentsRequest{
    Limit:  10,
    Status: xpay.PaymentStatusSucceeded,
})
```

## Working with Customers

### Create a Customer

```go
customer, err := client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
    Name:  "John Doe",
    Email: "john@example.com",
    Phone: "+231700000000",
})
```

### List Customers

```go
customers, err := client.Customers.List(ctx, &xpay.ListCustomersRequest{
    Limit: 10,
    Email: "john@example.com",
})
```

## Working with Webhooks

### Create a Webhook

```go
webhook, err := client.Webhooks.Create(ctx, &xpay.CreateWebhookRequest{
    URL:    "https://your-site.com/webhooks/xpay",
    Events: []string{"payment.succeeded", "payment.failed"},
})
```

### Verify Webhook Signatures

```go
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    payload, _ := io.ReadAll(r.Body)
    signature := r.Header.Get("X-XPay-Signature")
    secret := "your_webhook_secret"

    if !xpay.VerifyWebhookSignature(payload, signature, secret) {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // Process webhook...
    w.WriteHeader(http.StatusOK)
}
```

## Error Handling

```go
payment, err := client.Payments.Create(ctx, paymentRequest)
if err != nil {
    if xpayErr, ok := err.(*xpay.Error); ok {
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
// Convert to smallest unit (cents)
amount := decimal.NewFromFloat(25.99)
cents := xpay.ToSmallestUnit(amount, xpay.CurrencyUSD) // 2599

// Convert from smallest unit
dollars := xpay.FromSmallestUnit(cents, xpay.CurrencyUSD) // 25.99

// Format for display
formatted := xpay.FormatAmount(cents, xpay.CurrencyUSD, true) // "$25.99"
```

## Payment Methods

The SDK supports the following payment methods:

- `xpay.PaymentMethodStripe` - Credit/debit cards via Stripe
- `xpay.PaymentMethodMomoLiberia` - MTN Mobile Money Liberia
- `xpay.PaymentMethodXPayWallet` - X-Pay wallet transfers

## Testing

Run the test suite:

```bash
# Run unit tests
go test ./xpay

# Run the SDK test suite
go run ./test_sdk.go

# Run examples
go run ./examples/basic_payment.go
go run ./examples/customers.go
go run ./examples/webhooks.go
```

## Examples

Check out the `examples/` directory for complete working examples:

- `basic_payment.go` - Basic payment creation
- `customers.go` - Customer management
- `webhooks.go` - Webhook management

## Support

- Documentation: https://docs.xpay-bits.com
- Support: support@xpay-bits.com
- Issues: GitHub Issues