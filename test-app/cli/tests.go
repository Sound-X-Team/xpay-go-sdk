package cli

import (
	"bufio"
	"context"
	"fmt"

	"github.com/Sound-X-Team/xpay-go-sdk/xpay"
	"github.com/shopspring/decimal"
)

var (
	lastPaymentID  string
	lastCustomerID string
	lastWebhookID  string
)

// testConnectivity tests basic API connectivity
func (r *TestRunner) testConnectivity(ctx context.Context) error {
	ping, err := r.client.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	fmt.Printf("      API responded successfully (Success: %v)\n", ping.Success)
	return nil
}

// testCurrencyUtils tests currency utility functions
func (r *TestRunner) testCurrencyUtils(ctx context.Context) error {
	amount := decimal.NewFromFloat(25.99)
	cents := xpay.ToSmallestUnit(amount, xpay.CurrencyUSD)
	dollarsBack := xpay.FromSmallestUnit(cents, xpay.CurrencyUSD)
	formatted := xpay.FormatAmount(cents, xpay.CurrencyUSD, true)

	if cents != 2599 {
		return fmt.Errorf("expected 2599 cents, got %d", cents)
	}

	if !dollarsBack.Equal(amount) {
		return fmt.Errorf("currency conversion failed: %s != %s", dollarsBack.String(), amount.String())
	}

	if formatted != "$25.99" {
		return fmt.Errorf("expected '$25.99', got '%s'", formatted)
	}

	fmt.Printf("      Currency conversion: %s -> %d cents -> %s -> %s\n", 
		amount.String(), cents, dollarsBack.String(), formatted)
	return nil
}

// testStripePayment creates a Stripe payment
func (r *TestRunner) testStripePayment(ctx context.Context) error {
	payment, err := r.client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(29.99),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodStripe,
		Description:   "Go SDK Stripe test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			PaymentMethodTypes: []string{"card"},
		},
		Metadata: map[string]interface{}{
			"source": "go-sdk-test-app",
			"test":   true,
		},
	})

	if err != nil {
		return fmt.Errorf("stripe payment creation failed: %w", err)
	}

	lastPaymentID = payment.ID
	fmt.Printf("      Payment ID: %s\n", payment.ID)
	fmt.Printf("      Status: %s\n", payment.Status)
	fmt.Printf("      Amount: %s %s\n", payment.Amount.String(), payment.Currency)
	if payment.ClientSecret != "" {
		fmt.Printf("      Client Secret: %s...\n", payment.ClientSecret[:20])
	}

	return nil
}

// testMoMoPayment creates a Mobile Money payment
func (r *TestRunner) testMoMoPayment(ctx context.Context) error {
	payment, err := r.client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(50.00),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodMomoLiberia,
		Description:   "Go SDK MoMo test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			PhoneNumber: "+231700000000",
		},
		Metadata: map[string]interface{}{
			"source": "go-sdk-test-app",
			"test":   true,
		},
	})

	if err != nil {
		return fmt.Errorf("momo payment creation failed: %w", err)
	}

	fmt.Printf("      Payment ID: %s\n", payment.ID)
	fmt.Printf("      Status: %s\n", payment.Status)
	fmt.Printf("      Amount: %s %s\n", payment.Amount.String(), payment.Currency)
	if payment.ReferenceID != "" {
		fmt.Printf("      Reference ID: %s\n", payment.ReferenceID)
	}

	return nil
}

// testPaymentRetrieval retrieves a payment by ID
func (r *TestRunner) testPaymentRetrieval(ctx context.Context) error {
	if lastPaymentID == "" {
		return fmt.Errorf("no payment ID available (create a payment first)")
	}

	payment, err := r.client.Payments.Retrieve(ctx, lastPaymentID)
	if err != nil {
		return fmt.Errorf("payment retrieval failed: %w", err)
	}

	fmt.Printf("      Retrieved payment: %s\n", payment.ID)
	fmt.Printf("      Status: %s\n", payment.Status)
	fmt.Printf("      Amount: %s %s\n", payment.Amount.String(), payment.Currency)

	return nil
}

// testPaymentListing lists recent payments
func (r *TestRunner) testPaymentListing(ctx context.Context) error {
	payments, err := r.client.Payments.List(ctx, &xpay.ListPaymentsRequest{
		Limit: 5,
	})

	if err != nil {
		return fmt.Errorf("payment listing failed: %w", err)
	}

	fmt.Printf("      Found %d payments (total: %d)\n", len(payments.Items), payments.Total)
	for i, payment := range payments.Items {
		if i < 3 { // Show first 3
			fmt.Printf("        %d. %s - %s %s (%s)\n", 
				i+1, payment.ID, payment.Amount.String(), payment.Currency, payment.Status)
		}
	}

	return nil
}

// testCustomerManagement tests customer CRUD operations
func (r *TestRunner) testCustomerManagement(ctx context.Context) error {
	// Create customer
	customer, err := r.client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
		Name:  "Test Customer",
		Email: "test@example.com",
		Phone: "+231700000000",
		Metadata: map[string]interface{}{
			"source": "go-sdk-test-app",
			"test":   true,
		},
	})

	if err != nil {
		return fmt.Errorf("customer creation failed: %w", err)
	}

	lastCustomerID = customer.ID
	fmt.Printf("      Created customer: %s (%s)\n", customer.ID, customer.Name)

	// Retrieve customer
	retrievedCustomer, err := r.client.Customers.Retrieve(ctx, customer.ID)
	if err != nil {
		return fmt.Errorf("customer retrieval failed: %w", err)
	}

	fmt.Printf("      Retrieved customer: %s\n", retrievedCustomer.Name)

	// Update customer
	updatedCustomer, err := r.client.Customers.Update(ctx, customer.ID, &xpay.UpdateCustomerRequest{
		Name: "Updated Test Customer",
		Metadata: map[string]interface{}{
			"source":  "go-sdk-test-app",
			"test":    true,
			"updated": true,
		},
	})

	if err != nil {
		return fmt.Errorf("customer update failed: %w", err)
	}

	fmt.Printf("      Updated customer name: %s\n", updatedCustomer.Name)

	// List customers
	customers, err := r.client.Customers.List(ctx, &xpay.ListCustomersRequest{
		Limit: 3,
	})

	if err != nil {
		return fmt.Errorf("customer listing failed: %w", err)
	}

	fmt.Printf("      Listed %d customers (total: %d)\n", len(customers.Items), customers.Total)

	return nil
}

// testWebhookManagement tests webhook CRUD operations
func (r *TestRunner) testWebhookManagement(ctx context.Context) error {
	// Create webhook
	webhook, err := r.client.Webhooks.Create(ctx, &xpay.CreateWebhookRequest{
		URL: "https://example.com/webhook/xpay",
		Events: []string{
			"payment.succeeded",
			"payment.failed",
			"customer.created",
		},
		Description: "Go SDK test webhook",
	})

	if err != nil {
		return fmt.Errorf("webhook creation failed: %w", err)
	}

	lastWebhookID = webhook.ID
	fmt.Printf("      Created webhook: %s\n", webhook.ID)
	fmt.Printf("      URL: %s\n", webhook.URL)
	fmt.Printf("      Events: %v\n", webhook.Events)

	// List webhooks
	webhooks, err := r.client.Webhooks.List(ctx)
	if err != nil {
		return fmt.Errorf("webhook listing failed: %w", err)
	}

	fmt.Printf("      Listed %d webhooks\n", webhooks.Total)

	// Test signature verification
	payload := []byte(`{"event": "payment.succeeded", "data": {"id": "pay_123"}}`)
	signature := "test_signature"
	secret := webhook.Secret

	isValid := xpay.VerifyWebhookSignature(payload, signature, secret)
	fmt.Printf("      Signature verification: %v (expected false with test data)\n", isValid)

	return nil
}

// testErrorHandling tests error handling scenarios
func (r *TestRunner) testErrorHandling(ctx context.Context) error {
	// Try to retrieve a non-existent payment
	_, err := r.client.Payments.Retrieve(ctx, "pay_nonexistent_12345")
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			fmt.Printf("      Caught expected error: [%s] %s\n", xpayErr.Code, xpayErr.Message)
			return nil
		}
		return fmt.Errorf("expected XPayError, got %T", err)
	}

	return fmt.Errorf("expected error for non-existent payment")
}

// Interactive functions

func (r *TestRunner) interactivePayment(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("💳 Create a test payment")
	
	amount, err := getDecimalInput(reader, "Enter amount (e.g., 25.99): ")
	if err != nil {
		fmt.Printf("Error reading amount: %v\n", err)
		return
	}

	currency, err := getUserInput(reader, "Enter currency (USD, EUR, GBP, GHS): ")
	if err != nil {
		fmt.Printf("Error reading currency: %v\n", err)
		return
	}

	fmt.Println("Payment methods:")
	fmt.Println("  1. stripe")
	fmt.Println("  2. momo_liberia")
	fmt.Println("  3. xpay_wallet")

	methodChoice, err := getUserInput(reader, "Choose payment method (1-3): ")
	if err != nil {
		fmt.Printf("Error reading payment method: %v\n", err)
		return
	}

	var paymentMethod xpay.PaymentMethod
	var paymentMethodData *xpay.PaymentMethodData

	switch methodChoice {
	case "1", "stripe":
		paymentMethod = xpay.PaymentMethodStripe
		paymentMethodData = &xpay.PaymentMethodData{
			PaymentMethodTypes: []string{"card"},
		}
	case "2", "momo_liberia":
		paymentMethod = xpay.PaymentMethodMomoLiberia
		phone, _ := getUserInput(reader, "Enter phone number (e.g., +231700000000): ")
		paymentMethodData = &xpay.PaymentMethodData{
			PhoneNumber: phone,
		}
	case "3", "xpay_wallet":
		paymentMethod = xpay.PaymentMethodXPayWallet
		walletID, _ := getUserInput(reader, "Enter wallet ID: ")
		pin, _ := getUserInput(reader, "Enter PIN: ")
		paymentMethodData = &xpay.PaymentMethodData{
			WalletID: walletID,
			PIN:      pin,
		}
	default:
		fmt.Println("Invalid choice")
		return
	}

	description, _ := getUserInput(reader, "Enter description (optional): ")

	payment, err := r.client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:            amount,
		Currency:          xpay.Currency(currency),
		PaymentMethod:     paymentMethod,
		Description:       description,
		PaymentMethodData: paymentMethodData,
		Metadata: map[string]interface{}{
			"source":      "interactive-test",
			"created_via": "go-sdk-test-app",
		},
	})

	if err != nil {
		fmt.Printf("❌ Payment creation failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Payment created successfully!\n")
	fmt.Printf("   ID: %s\n", payment.ID)
	fmt.Printf("   Status: %s\n", payment.Status)
	fmt.Printf("   Amount: %s %s\n", payment.Amount.String(), payment.Currency)

	lastPaymentID = payment.ID
}

func (r *TestRunner) interactiveCustomer(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("👤 Create a test customer")

	name, err := getUserInput(reader, "Enter customer name: ")
	if err != nil {
		fmt.Printf("Error reading name: %v\n", err)
		return
	}

	email, _ := getUserInput(reader, "Enter email (optional): ")
	phone, _ := getUserInput(reader, "Enter phone (optional): ")

	customer, err := r.client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
		Name:  name,
		Email: email,
		Phone: phone,
		Metadata: map[string]interface{}{
			"source":      "interactive-test",
			"created_via": "go-sdk-test-app",
		},
	})

	if err != nil {
		fmt.Printf("❌ Customer creation failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Customer created successfully!\n")
	fmt.Printf("   ID: %s\n", customer.ID)
	fmt.Printf("   Name: %s\n", customer.Name)
	fmt.Printf("   Email: %s\n", customer.Email)

	lastCustomerID = customer.ID
}

func (r *TestRunner) interactiveWebhook(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("🔗 Create a test webhook")

	url, err := getUserInput(reader, "Enter webhook URL: ")
	if err != nil {
		fmt.Printf("Error reading URL: %v\n", err)
		return
	}

	webhook, err := r.client.Webhooks.Create(ctx, &xpay.CreateWebhookRequest{
		URL: url,
		Events: []string{
			"payment.succeeded",
			"payment.failed",
			"customer.created",
		},
		Description: "Interactive test webhook",
	})

	if err != nil {
		fmt.Printf("❌ Webhook creation failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Webhook created successfully!\n")
	fmt.Printf("   ID: %s\n", webhook.ID)
	fmt.Printf("   URL: %s\n", webhook.URL)
	fmt.Printf("   Events: %v\n", webhook.Events)
	fmt.Printf("   Secret: %s\n", webhook.Secret)

	lastWebhookID = webhook.ID
}

func (r *TestRunner) interactiveListPayments(ctx context.Context) {
	fmt.Println("📋 Listing recent payments")

	payments, err := r.client.Payments.List(ctx, &xpay.ListPaymentsRequest{
		Limit: 10,
	})

	if err != nil {
		fmt.Printf("❌ Failed to list payments: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d payments (total: %d)\n", len(payments.Items), payments.Total)
	for i, payment := range payments.Items {
		fmt.Printf("   %d. %s - %s %s - %s\n", 
			i+1, payment.ID, payment.Amount.String(), payment.Currency, payment.Status)
	}
}