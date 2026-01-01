package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/test-app/config"
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
	"github.com/shopspring/decimal"
)

// TestRunner handles CLI test execution
type TestRunner struct {
	client *xpay.Client
	cfg    *config.Config
}

// NewTestRunner creates a new test runner
func NewTestRunner(cfg *config.Config) *TestRunner {
	client := xpay.NewClient(cfg.XPay)
	return &TestRunner{
		client: client,
		cfg:    cfg,
	}
}

// RunTests runs the comprehensive test suite
func RunTests(cfg *config.Config) error {
	runner := NewTestRunner(cfg)
	
	if cfg.Debug {
		cfg.Print()
	}

	fmt.Println("🧪 Running X-Pay Go SDK Test Suite")
	fmt.Println("==================================")
	
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TestTimeout)
	defer cancel()

	tests := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"API Connectivity", runner.testConnectivity},
		{"Currency Utilities", runner.testCurrencyUtils},
		{"Payment Creation (Stripe)", runner.testStripePayment},
		{"Payment Creation (MoMo)", runner.testMoMoPayment},
		{"Payment Retrieval", runner.testPaymentRetrieval},
		{"Payment Listing", runner.testPaymentListing},
		{"Customer Management", runner.testCustomerManagement},
		{"Webhook Management", runner.testWebhookManagement},
		{"Error Handling", runner.testErrorHandling},
	}

	passed := 0
	total := len(tests)

	for i, test := range tests {
		fmt.Printf("\n📋 Test %d/%d: %s\n", i+1, total, test.name)
		
		if err := test.fn(ctx); err != nil {
			fmt.Printf("   ❌ FAILED: %v\n", err)
		} else {
			fmt.Printf("   ✅ PASSED\n")
			passed++
		}
	}

	fmt.Printf("\n📊 Test Results\n")
	fmt.Println("===============")
	fmt.Printf("Passed: %d/%d tests\n", passed, total)
	fmt.Printf("Success Rate: %.1f%%\n", float64(passed)/float64(total)*100)

	if passed == total {
		fmt.Println("🎉 All tests passed!")
		return nil
	} else {
		return fmt.Errorf("%d tests failed", total-passed)
	}
}

// RunPaymentTests runs only payment-related tests
func RunPaymentTests(cfg *config.Config) error {
	runner := NewTestRunner(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TestTimeout)
	defer cancel()

	fmt.Println("💳 Testing Payment Operations")
	fmt.Println("=============================")

	tests := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Stripe Payment", runner.testStripePayment},
		{"Mobile Money Payment", runner.testMoMoPayment},
		{"Payment Retrieval", runner.testPaymentRetrieval},
		{"Payment Listing", runner.testPaymentListing},
	}

	for _, test := range tests {
		fmt.Printf("\n🔍 %s: ", test.name)
		if err := test.fn(ctx); err != nil {
			fmt.Printf("❌ FAILED - %v\n", err)
			return err
		}
		fmt.Printf("✅ PASSED\n")
	}

	fmt.Println("\n🎉 All payment tests passed!")
	return nil
}

// RunCustomerTests runs only customer-related tests
func RunCustomerTests(cfg *config.Config) error {
	runner := NewTestRunner(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TestTimeout)
	defer cancel()

	fmt.Println("👥 Testing Customer Operations")
	fmt.Println("==============================")

	if err := runner.testCustomerManagement(ctx); err != nil {
		fmt.Printf("❌ Customer tests failed: %v\n", err)
		return err
	}

	fmt.Println("✅ All customer tests passed!")
	return nil
}

// RunWebhookTests runs only webhook-related tests
func RunWebhookTests(cfg *config.Config) error {
	runner := NewTestRunner(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TestTimeout)
	defer cancel()

	fmt.Println("🔗 Testing Webhook Operations")
	fmt.Println("=============================")

	if err := runner.testWebhookManagement(ctx); err != nil {
		fmt.Printf("❌ Webhook tests failed: %v\n", err)
		return err
	}

	fmt.Println("✅ All webhook tests passed!")
	return nil
}

// PingAPI tests basic API connectivity
func PingAPI(cfg *config.Config) error {
	runner := NewTestRunner(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Testing API Connectivity")
	fmt.Println("===========================")

	if err := runner.testConnectivity(ctx); err != nil {
		fmt.Printf("❌ API connectivity failed: %v\n", err)
		return err
	}

	fmt.Println("✅ API is reachable!")
	return nil
}

// RunInteractive starts an interactive test session
func RunInteractive(cfg *config.Config) error {
	runner := NewTestRunner(cfg)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🎮 Interactive X-Pay SDK Test Session")
	fmt.Println("====================================")
	fmt.Println("Available commands:")
	fmt.Println("  1. ping      - Test API connectivity")
	fmt.Println("  2. payment   - Create a test payment")
	fmt.Println("  3. customer  - Create a test customer")
	fmt.Println("  4. webhook   - Create a test webhook")
	fmt.Println("  5. list      - List recent payments")
	fmt.Println("  6. help      - Show this help")
	fmt.Println("  7. quit      - Exit interactive mode")
	fmt.Println()

	for {
		fmt.Print("xpay> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		command := strings.TrimSpace(input)
		ctx := context.Background()

		switch command {
		case "ping", "1":
			runner.testConnectivity(ctx)

		case "payment", "2":
			runner.interactivePayment(ctx, reader)

		case "customer", "3":
			runner.interactiveCustomer(ctx, reader)

		case "webhook", "4":
			runner.interactiveWebhook(ctx, reader)

		case "list", "5":
			runner.interactiveListPayments(ctx)

		case "help", "6":
			fmt.Println("Available commands: ping, payment, customer, webhook, list, help, quit")

		case "quit", "exit", "7":
			fmt.Println("👋 Goodbye!")
			return nil

		case "":
			continue

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}

		fmt.Println()
	}
}

// Helper function to get user input
func getUserInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// Helper function to get decimal input
func getDecimalInput(reader *bufio.Reader, prompt string) (decimal.Decimal, error) {
	input, err := getUserInput(reader, prompt)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(input)
}