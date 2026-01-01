package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
	"github.com/shopspring/decimal"
)

func main() {
	// Initialize client with working test credentials
	client := xpay.NewClient(&xpay.Config{
		APIKey:      "sk_sandbox_7c845adf-f658-4f29-9857-7e8a8708",
		MerchantID:  "548d8033-fbe9-411b-991f-f159cdee7745",
		Environment: xpay.EnvironmentSandbox,
		BaseURL:     "http://localhost:8000",
	})

	ctx := context.Background()

	// Test API connectivity
	fmt.Println("Testing API connectivity...")
	ping, err := client.Ping(ctx)
	if err != nil {
		log.Printf("Ping failed: %v", err)
	} else {
		fmt.Printf("✅ API is reachable - Success: %v, Timestamp: %s\n", ping.Success, ping.Timestamp)
	}

	// Create a Stripe payment
	fmt.Println("\nCreating Stripe payment...")
	stripePayment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(25.99),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodStripe,
		Description:   "Go SDK Stripe test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			PaymentMethodTypes: []string{"card"},
		},
		Metadata: map[string]interface{}{
			"source": "go-sdk-example",
			"test":   true,
		},
	})

	if err != nil {
		log.Printf("Stripe payment creation failed: %v", err)
	} else {
		fmt.Printf("✅ Stripe payment created: %s\n", stripePayment.ID)
		fmt.Printf("   Status: %s\n", stripePayment.Status)
		fmt.Printf("   Amount: %s %s\n", stripePayment.Amount.String(), stripePayment.Currency)
		if stripePayment.ClientSecret != "" {
			fmt.Printf("   Client Secret: %s\n", stripePayment.ClientSecret)
		}
	}

	// Create a Mobile Money payment
	fmt.Println("\nCreating Mobile Money payment...")
	momoPayment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(50.00),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodMomoLiberia,
		Description:   "Go SDK MoMo test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			PhoneNumber: "+231700000000",
		},
		Metadata: map[string]interface{}{
			"source": "go-sdk-example",
			"test":   true,
		},
	})

	if err != nil {
		log.Printf("MoMo payment creation failed: %v", err)
	} else {
		fmt.Printf("✅ MoMo payment created: %s\n", momoPayment.ID)
		fmt.Printf("   Status: %s\n", momoPayment.Status)
		fmt.Printf("   Amount: %s %s\n", momoPayment.Amount.String(), momoPayment.Currency)
		if momoPayment.ReferenceID != "" {
			fmt.Printf("   Reference ID: %s\n", momoPayment.ReferenceID)
		}
		if momoPayment.Instructions != "" {
			fmt.Printf("   Instructions: %s\n", momoPayment.Instructions)
		}
	}

	// Create a X-Pay Wallet payment
	fmt.Println("\nCreating X-Pay Wallet payment...")
	walletPayment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(15.50),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodXPayWallet,
		Description:   "Go SDK Wallet test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			WalletID: "wallet_123",
			PIN:      "1234",
		},
		Metadata: map[string]interface{}{
			"source": "go-sdk-example",
			"test":   true,
		},
	})

	if err != nil {
		log.Printf("Wallet payment creation failed: %v", err)
	} else {
		fmt.Printf("✅ Wallet payment created: %s\n", walletPayment.ID)
		fmt.Printf("   Status: %s\n", walletPayment.Status)
		fmt.Printf("   Amount: %s %s\n", walletPayment.Amount.String(), walletPayment.Currency)
	}

	// List recent payments
	fmt.Println("\nListing recent payments...")
	payments, err := client.Payments.List(ctx, &xpay.ListPaymentsRequest{
		Limit: 5,
	})

	if err != nil {
		log.Printf("Failed to list payments: %v", err)
	} else {
		fmt.Printf("✅ Found %d payments (showing %d):\n", payments.Total, len(payments.Items))
		for i, payment := range payments.Items {
			fmt.Printf("   %d. %s - %s %s - %s\n", 
				i+1, payment.ID, payment.Amount.String(), payment.Currency, payment.Status)
		}
	}

	// Test currency utilities
	fmt.Println("\nTesting currency utilities...")
	amount := decimal.NewFromFloat(25.99)
	cents := xpay.ToSmallestUnit(amount, xpay.CurrencyUSD)
	dollarsBack := xpay.FromSmallestUnit(cents, xpay.CurrencyUSD)
	formatted := xpay.FormatAmount(cents, xpay.CurrencyUSD, true)

	fmt.Printf("Original amount: %s\n", amount.String())
	fmt.Printf("In cents: %d\n", cents)
	fmt.Printf("Back to dollars: %s\n", dollarsBack.String())
	fmt.Printf("Formatted: %s\n", formatted)

	// Test payment method validation
	fmt.Println("\nTesting supported currencies...")
	supportedCurrencies := client.Payments.GetSupportedCurrencies(xpay.PaymentMethodStripe)
	fmt.Printf("Stripe supported currencies: %v\n", supportedCurrencies)

	supportedCurrencies = client.Payments.GetSupportedCurrencies(xpay.PaymentMethodMomoLiberia)
	fmt.Printf("MoMo Liberia supported currencies: %v\n", supportedCurrencies)

	fmt.Println("\n🎉 Go SDK test completed successfully!")
}