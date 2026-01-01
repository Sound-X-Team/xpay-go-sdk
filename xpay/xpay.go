// Package xpay provides the official Go SDK for X-Pay payment processing platform.
// Accept payments from multiple providers including Stripe, Mobile Money, and X-Pay Wallets with a unified API.
//
// Example usage:
//
//	client := xpay.NewClient(&xpay.Config{
//		APIKey:      "sk_sandbox_your_secret_key_here",
//		MerchantID:  "your_merchant_id_here",
//		Environment: xpay.EnvironmentSandbox,
//	})
//
//	payment, err := client.Payments.Create(context.Background(), &xpay.PaymentRequest{
//		Amount:        decimal.NewFromFloat(29.99),
//		Currency:      xpay.CurrencyUSD,
//		PaymentMethod: xpay.PaymentMethodStripe,
//		Description:   "Test payment",
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Payment created: %s\n", payment.ID)
package xpay