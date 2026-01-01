package main

import (
	"context"
	"fmt"
	"os"
	"time"
	
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("🧪 X-Pay Go SDK Test Suite")
	fmt.Println("==========================")

	// Initialize client with test credentials
	client := xpay.NewClient(&xpay.Config{
		APIKey:      "sk_sandbox_7c845adf-f658-4f29-9857-7e8a8708",
		MerchantID:  "548d8033-fbe9-411b-991f-f159cdee7745",
		Environment: xpay.EnvironmentSandbox,
		BaseURL:     "http://localhost:8000",
		Timeout:     30 * time.Second,
	})

	ctx := context.Background()
	
	// Track test results
	passedTests := 0
	totalTests := 0

	// Test 1: API Connectivity
	totalTests++
	fmt.Printf("\n📡 Test %d: API Connectivity\n", totalTests)
	fmt.Println("   Testing ping endpoint...")
	
	ping, err := client.Ping(ctx)
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: API is reachable (Success: %v)\n", ping.Success)
		passedTests++
	}

	// Test 2: Currency Utilities
	totalTests++
	fmt.Printf("\n💰 Test %d: Currency Utilities\n", totalTests)
	
	amount := decimal.NewFromFloat(25.99)
	cents := xpay.ToSmallestUnit(amount, xpay.CurrencyUSD)
	dollarsBack := xpay.FromSmallestUnit(cents, xpay.CurrencyUSD)
	formatted := xpay.FormatAmount(cents, xpay.CurrencyUSD, true)
	
	if cents == 2599 && dollarsBack.Equal(amount) && formatted == "$25.99" {
		fmt.Printf("   ✅ PASS: Currency conversion works correctly\n")
		fmt.Printf("      %s -> %d cents -> %s -> %s\n", amount.String(), cents, dollarsBack.String(), formatted)
		passedTests++
	} else {
		fmt.Printf("   ❌ FAIL: Currency conversion failed\n")
		fmt.Printf("      Expected: 2599 cents, $25.99\n")
		fmt.Printf("      Got: %d cents, %s\n", cents, formatted)
	}

	// Test 3: Payment Method Validation
	totalTests++
	fmt.Printf("\n🔍 Test %d: Payment Method Validation\n", totalTests)
	
	supportedCurrencies := client.Payments.GetSupportedCurrencies(xpay.PaymentMethodStripe)
	expectedStripe := []xpay.Currency{xpay.CurrencyUSD, xpay.CurrencyEUR, xpay.CurrencyGBP, xpay.CurrencyGHS}
	
	if len(supportedCurrencies) == len(expectedStripe) {
		fmt.Printf("   ✅ PASS: Stripe supports expected currencies: %v\n", supportedCurrencies)
		passedTests++
	} else {
		fmt.Printf("   ❌ FAIL: Expected %v, got %v\n", expectedStripe, supportedCurrencies)
	}

	// Test 4: Stripe Payment Creation
	totalTests++
	fmt.Printf("\n💳 Test %d: Stripe Payment Creation\n", totalTests)
	
	stripePayment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(29.99),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodStripe,
		Description:   "Go SDK test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			PaymentMethodTypes: []string{"card"},
		},
		Metadata: map[string]interface{}{
			"test": true,
			"sdk":  "go",
		},
	})
	
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: Stripe payment created successfully\n")
		fmt.Printf("      ID: %s\n", stripePayment.ID)
		fmt.Printf("      Status: %s\n", stripePayment.Status)
		fmt.Printf("      Amount: %s %s\n", stripePayment.Amount.String(), stripePayment.Currency)
		passedTests++
	}

	// Test 5: Mobile Money Payment Creation
	totalTests++
	fmt.Printf("\n📱 Test %d: Mobile Money Payment Creation\n", totalTests)
	
	momoPayment, err := client.Payments.Create(ctx, &xpay.PaymentRequest{
		Amount:        decimal.NewFromFloat(50.00),
		Currency:      xpay.CurrencyUSD,
		PaymentMethod: xpay.PaymentMethodMomoLiberia,
		Description:   "Go SDK MoMo test payment",
		PaymentMethodData: &xpay.PaymentMethodData{
			PhoneNumber: "+231700000000",
		},
		Metadata: map[string]interface{}{
			"test": true,
			"sdk":  "go",
		},
	})
	
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: MoMo payment created successfully\n")
		fmt.Printf("      ID: %s\n", momoPayment.ID)
		fmt.Printf("      Status: %s\n", momoPayment.Status)
		fmt.Printf("      Amount: %s %s\n", momoPayment.Amount.String(), momoPayment.Currency)
		passedTests++
	}

	// Test 6: Payment Retrieval (if we have a payment ID)
	var testPaymentID string
	if stripePayment != nil {
		testPaymentID = stripePayment.ID
	} else if momoPayment != nil {
		testPaymentID = momoPayment.ID
	}

	if testPaymentID != "" {
		totalTests++
		fmt.Printf("\n🔍 Test %d: Payment Retrieval\n", totalTests)
		
		retrievedPayment, err := client.Payments.Retrieve(ctx, testPaymentID)
		if err != nil {
			fmt.Printf("   ❌ FAIL: %v\n", err)
		} else {
			fmt.Printf("   ✅ PASS: Payment retrieved successfully\n")
			fmt.Printf("      ID: %s\n", retrievedPayment.ID)
			fmt.Printf("      Status: %s\n", retrievedPayment.Status)
			passedTests++
		}
	}

	// Test 7: Payment Listing
	totalTests++
	fmt.Printf("\n📋 Test %d: Payment Listing\n", totalTests)
	
	payments, err := client.Payments.List(ctx, &xpay.ListPaymentsRequest{
		Limit: 5,
	})
	
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: Listed %d payments (total: %d)\n", len(payments.Items), payments.Total)
		for i, payment := range payments.Items {
			if i < 3 { // Show first 3
				fmt.Printf("      %d. %s - %s %s\n", i+1, payment.ID, payment.Amount.String(), payment.Currency)
			}
		}
		passedTests++
	}

	// Test 8: Customer Creation
	totalTests++
	fmt.Printf("\n👤 Test %d: Customer Creation\n", totalTests)
	
	customer, err := client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
		Name:  "Test Customer",
		Email: "test@example.com",
		Phone: "+231700000000",
		Metadata: map[string]interface{}{
			"test": true,
			"sdk":  "go",
		},
	})
	
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: Customer created successfully\n")
		fmt.Printf("      ID: %s\n", customer.ID)
		fmt.Printf("      Name: %s\n", customer.Name)
		fmt.Printf("      Email: %s\n", customer.Email)
		passedTests++
	}

	// Test 9: Customer Listing
	totalTests++
	fmt.Printf("\n👥 Test %d: Customer Listing\n", totalTests)
	
	customers, err := client.Customers.List(ctx, &xpay.ListCustomersRequest{
		Limit: 3,
	})
	
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: Listed %d customers (total: %d)\n", len(customers.Items), customers.Total)
		for i, c := range customers.Items {
			fmt.Printf("      %d. %s - %s\n", i+1, c.ID, c.Name)
		}
		passedTests++
	}

	// Test 10: Error Handling
	totalTests++
	fmt.Printf("\n❗ Test %d: Error Handling\n", totalTests)
	
	_, err = client.Payments.Retrieve(ctx, "pay_nonexistent_payment_id")
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			fmt.Printf("   ✅ PASS: Error handling works correctly\n")
			fmt.Printf("      Error Code: %s\n", xpayErr.Code)
			fmt.Printf("      Error Message: %s\n", xpayErr.Message)
			passedTests++
		} else {
			fmt.Printf("   ❌ FAIL: Expected XPayError, got %T\n", err)
		}
	} else {
		fmt.Printf("   ❌ FAIL: Expected error for nonexistent payment\n")
	}

	// Test Summary
	fmt.Printf("\n📊 Test Results Summary\n")
	fmt.Println("=======================")
	fmt.Printf("Passed: %d/%d tests\n", passedTests, totalTests)
	fmt.Printf("Success Rate: %.1f%%\n", float64(passedTests)/float64(totalTests)*100)

	if passedTests == totalTests {
		fmt.Println("🎉 All tests passed! SDK is working correctly.")
		os.Exit(0)
	} else {
		fmt.Printf("⚠️  %d tests failed. Please check the implementation.\n", totalTests-passedTests)
		os.Exit(1)
	}
}