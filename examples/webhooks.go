package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	
	"github.com/Sound-X-Team/xpay-go-sdk/xpay"
)

func main() {
	// Initialize client
	client := xpay.NewClient(&xpay.Config{
		APIKey:      "sk_sandbox_7c845adf-f658-4f29-9857-7e8a8708",
		MerchantID:  "548d8033-fbe9-411b-991f-f159cdee7745",
		Environment: xpay.EnvironmentSandbox,
		BaseURL:     "http://localhost:8000",
	})

	ctx := context.Background()

	// Create a webhook endpoint
	fmt.Println("Creating webhook endpoint...")
	webhook, err := client.Webhooks.Create(ctx, &xpay.CreateWebhookRequest{
		URL: "https://your-site.com/webhooks/xpay",
		Events: []string{
			"payment.succeeded",
			"payment.failed",
			"payment.pending",
			"customer.created",
		},
		Description: "Go SDK webhook example",
	})

	if err != nil {
		log.Printf("Webhook creation failed: %v", err)
		return
	}

	fmt.Printf("✅ Webhook created: %s\n", webhook.ID)
	fmt.Printf("   URL: %s\n", webhook.URL)
	fmt.Printf("   Events: %v\n", webhook.Events)
	fmt.Printf("   Environment: %s\n", webhook.Environment)
	fmt.Printf("   Active: %v\n", webhook.IsActive)
	fmt.Printf("   Secret: %s\n", webhook.Secret)

	// Retrieve the webhook
	fmt.Println("\nRetrieving webhook...")
	retrievedWebhook, err := client.Webhooks.Retrieve(ctx, webhook.ID)
	if err != nil {
		log.Printf("Webhook retrieval failed: %v", err)
	} else {
		fmt.Printf("✅ Webhook retrieved: %s - %s\n", retrievedWebhook.ID, retrievedWebhook.URL)
	}

	// Update the webhook
	fmt.Println("\nUpdating webhook...")
	isActive := false
	updatedWebhook, err := client.Webhooks.Update(ctx, webhook.ID, &xpay.UpdateWebhookRequest{
		Events: []string{
			"payment.succeeded",
			"payment.failed",
			"payment.cancelled",
			"customer.created",
			"customer.updated",
		},
		IsActive:    &isActive,
		Description: "Updated Go SDK webhook example",
	})

	if err != nil {
		log.Printf("Webhook update failed: %v", err)
	} else {
		fmt.Printf("✅ Webhook updated: %s\n", updatedWebhook.ID)
		fmt.Printf("   New Events: %v\n", updatedWebhook.Events)
		fmt.Printf("   Active: %v\n", updatedWebhook.IsActive)
	}

	// List all webhooks
	fmt.Println("\nListing webhooks...")
	webhooks, err := client.Webhooks.List(ctx)
	if err != nil {
		log.Printf("Failed to list webhooks: %v", err)
	} else {
		fmt.Printf("✅ Found %d webhooks:\n", webhooks.Total)
		for i, w := range webhooks.Items {
			fmt.Printf("   %d. %s - %s (Active: %v)\n", i+1, w.ID, w.URL, w.IsActive)
			fmt.Printf("      Events: %v\n", w.Events)
		}
	}

	// Test the webhook
	fmt.Println("\nTesting webhook...")
	testResult, err := client.Webhooks.Test(ctx, webhook.ID)
	if err != nil {
		log.Printf("Webhook test failed: %v", err)
	} else {
		fmt.Printf("✅ Webhook test result: Success=%v, Message=%s\n", 
			testResult.Success, testResult.Message)
	}

	// Demonstrate webhook signature verification
	fmt.Println("\nDemonstrating webhook signature verification...")
	
	// Simulate a webhook payload
	webhookPayload := `{
		"event": "payment.succeeded",
		"data": {
			"id": "pay_123456789",
			"status": "succeeded",
			"amount": "25.99",
			"currency": "USD"
		},
		"timestamp": "2024-01-20T10:30:00Z"
	}`
	
	// This would be the signature from the X-Pay-Signature header
	webhookSecret := webhook.Secret
	
	// In a real webhook handler, you'd get the signature from the request headers
	// For this example, we'll demonstrate how to verify it
	fmt.Printf("Webhook secret: %s\n", webhookSecret)
	fmt.Printf("Payload: %s\n", webhookPayload)
	
	// In your actual webhook handler, you would do something like:
	fmt.Println("\nExample webhook handler:")
	fmt.Println("```go")
	fmt.Println("func webhookHandler(w http.ResponseWriter, r *http.Request) {")
	fmt.Println("    payload, _ := io.ReadAll(r.Body)")
	fmt.Println("    signature := r.Header.Get(\"X-XPay-Signature\")")
	fmt.Println("    secret := \"your_webhook_secret\"")
	fmt.Println("    ")
	fmt.Println("    if !xpay.VerifyWebhookSignature(payload, signature, secret) {")
	fmt.Println("        http.Error(w, \"Invalid signature\", http.StatusUnauthorized)")
	fmt.Println("        return")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // Process the webhook...")
	fmt.Println("    w.WriteHeader(http.StatusOK)")
	fmt.Println("}")
	fmt.Println("```")

	// Create another webhook for testing
	fmt.Println("\nCreating a second webhook endpoint...")
	webhook2, err := client.Webhooks.Create(ctx, &xpay.CreateWebhookRequest{
		URL: "https://another-site.com/webhooks/xpay",
		Events: []string{
			"payment.succeeded",
			"refund.created",
		},
		Description: "Second webhook for testing",
	})

	if err != nil {
		log.Printf("Second webhook creation failed: %v", err)
	} else {
		fmt.Printf("✅ Second webhook created: %s\n", webhook2.ID)
	}

	// List webhooks again
	fmt.Println("\nListing all webhooks again...")
	allWebhooks, err := client.Webhooks.List(ctx)
	if err != nil {
		log.Printf("Failed to list all webhooks: %v", err)
	} else {
		fmt.Printf("✅ Total webhooks: %d\n", allWebhooks.Total)
		for i, w := range allWebhooks.Items {
			fmt.Printf("   %d. %s - %s\n", i+1, w.ID, w.URL)
		}
	}

	fmt.Println("\n🎉 Webhook management example completed successfully!")
}

// Example webhook handler function (for reference)
func exampleWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Get the signature from headers
	signature := r.Header.Get("X-XPay-Signature")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusBadRequest)
		return
	}

	// Your webhook secret (store this securely)
	webhookSecret := "whsec_your_webhook_secret_here"

	// Verify the signature
	if !xpay.VerifyWebhookSignature(payload, signature, webhookSecret) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse the webhook payload
	var webhookData map[string]interface{}
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Handle different event types
	eventType, ok := webhookData["event"].(string)
	if !ok {
		http.Error(w, "Missing event type", http.StatusBadRequest)
		return
	}

	switch eventType {
	case "payment.succeeded":
		// Handle successful payment
		fmt.Println("Payment succeeded!")
		
	case "payment.failed":
		// Handle failed payment
		fmt.Println("Payment failed!")
		
	case "customer.created":
		// Handle new customer
		fmt.Println("New customer created!")
		
	default:
		fmt.Printf("Unhandled event type: %s\n", eventType)
	}

	// Respond with 200 OK to acknowledge receipt
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}