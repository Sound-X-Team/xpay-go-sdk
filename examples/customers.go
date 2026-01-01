package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
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

	// Create a customer
	fmt.Println("Creating a customer...")
	customer, err := client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "+231700000000",
		Metadata: map[string]interface{}{
			"source":      "go-sdk-example",
			"signup_date": "2024-01-15",
			"vip":         true,
		},
	})

	if err != nil {
		log.Printf("Customer creation failed: %v", err)
		return
	}

	fmt.Printf("✅ Customer created: %s\n", customer.ID)
	fmt.Printf("   Name: %s\n", customer.Name)
	fmt.Printf("   Email: %s\n", customer.Email)
	fmt.Printf("   Phone: %s\n", customer.Phone)
	fmt.Printf("   Created: %s\n", customer.CreatedAt.Format("2006-01-02 15:04:05"))

	// Retrieve the customer
	fmt.Println("\nRetrieving customer...")
	retrievedCustomer, err := client.Customers.Retrieve(ctx, customer.ID)
	if err != nil {
		log.Printf("Customer retrieval failed: %v", err)
	} else {
		fmt.Printf("✅ Customer retrieved: %s - %s\n", retrievedCustomer.ID, retrievedCustomer.Name)
	}

	// Update the customer
	fmt.Println("\nUpdating customer...")
	updatedCustomer, err := client.Customers.Update(ctx, customer.ID, &xpay.UpdateCustomerRequest{
		Name:  "John Smith",
		Phone: "+231700000001",
		Metadata: map[string]interface{}{
			"source":         "go-sdk-example",
			"signup_date":    "2024-01-15",
			"vip":            true,
			"updated":        true,
			"last_modified":  "2024-01-20",
		},
	})

	if err != nil {
		log.Printf("Customer update failed: %v", err)
	} else {
		fmt.Printf("✅ Customer updated: %s\n", updatedCustomer.ID)
		fmt.Printf("   New Name: %s\n", updatedCustomer.Name)
		fmt.Printf("   New Phone: %s\n", updatedCustomer.Phone)
	}

	// List customers
	fmt.Println("\nListing customers...")
	customers, err := client.Customers.List(ctx, &xpay.ListCustomersRequest{
		Limit: 10,
	})

	if err != nil {
		log.Printf("Failed to list customers: %v", err)
	} else {
		fmt.Printf("✅ Found %d customers (showing %d):\n", customers.Total, len(customers.Items))
		for i, c := range customers.Items {
			fmt.Printf("   %d. %s - %s (%s)\n", i+1, c.ID, c.Name, c.Email)
		}
	}

	// Search for customers by email
	fmt.Println("\nSearching customers by email...")
	emailCustomers, err := client.Customers.List(ctx, &xpay.ListCustomersRequest{
		Email: "john.doe@example.com",
		Limit: 5,
	})

	if err != nil {
		log.Printf("Failed to search customers: %v", err)
	} else {
		fmt.Printf("✅ Found %d customers with email 'john.doe@example.com':\n", len(emailCustomers.Items))
		for i, c := range emailCustomers.Items {
			fmt.Printf("   %d. %s - %s\n", i+1, c.ID, c.Name)
		}
	}

	// Create another customer to demonstrate listing
	fmt.Println("\nCreating another customer...")
	customer2, err := client.Customers.Create(ctx, &xpay.CreateCustomerRequest{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
		Phone: "+231700000002",
		Metadata: map[string]interface{}{
			"source": "go-sdk-example",
			"tier":   "premium",
		},
	})

	if err != nil {
		log.Printf("Second customer creation failed: %v", err)
	} else {
		fmt.Printf("✅ Second customer created: %s - %s\n", customer2.ID, customer2.Name)
	}

	// List customers again to show both
	fmt.Println("\nListing all customers again...")
	allCustomers, err := client.Customers.List(ctx, &xpay.ListCustomersRequest{
		Limit: 20,
	})

	if err != nil {
		log.Printf("Failed to list all customers: %v", err)
	} else {
		fmt.Printf("✅ Total customers: %d\n", allCustomers.Total)
		fmt.Println("Recent customers:")
		for i, c := range allCustomers.Items {
			if i < 5 { // Show first 5
				fmt.Printf("   %d. %s - %s (%s)\n", i+1, c.ID, c.Name, c.Email)
			}
		}
	}

	fmt.Println("\n🎉 Customer management example completed successfully!")
}