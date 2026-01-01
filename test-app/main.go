package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/test-app/cli"
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/test-app/config"
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/test-app/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get command from args
	command := "server"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	fmt.Println("🚀 X-Pay Go SDK Test Application")
	fmt.Println("================================")

	switch command {
	case "server":
		fmt.Println("Starting HTTP server mode...")
		if err := server.Run(cfg); err != nil {
			log.Fatalf("Server failed: %v", err)
		}

	case "test":
		fmt.Println("Running comprehensive test suite...")
		if err := cli.RunTests(cfg); err != nil {
			log.Fatalf("Tests failed: %v", err)
		}

	case "test-payments":
		fmt.Println("Running payment tests...")
		if err := cli.RunPaymentTests(cfg); err != nil {
			log.Fatalf("Payment tests failed: %v", err)
		}

	case "test-customers":
		fmt.Println("Running customer tests...")
		if err := cli.RunCustomerTests(cfg); err != nil {
			log.Fatalf("Customer tests failed: %v", err)
		}

	case "test-webhooks":
		fmt.Println("Running webhook tests...")
		if err := cli.RunWebhookTests(cfg); err != nil {
			log.Fatalf("Webhook tests failed: %v", err)
		}

	case "ping":
		fmt.Println("Testing API connectivity...")
		if err := cli.PingAPI(cfg); err != nil {
			log.Fatalf("Ping failed: %v", err)
		}

	case "interactive":
		fmt.Println("Starting interactive mode...")
		if err := cli.RunInteractive(cfg); err != nil {
			log.Fatalf("Interactive mode failed: %v", err)
		}

	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  server           - Start HTTP server with REST API endpoints")
	fmt.Println("  test             - Run comprehensive test suite")
	fmt.Println("  test-payments    - Run payment tests only")
	fmt.Println("  test-customers   - Run customer tests only")
	fmt.Println("  test-webhooks    - Run webhook tests only")
	fmt.Println("  ping             - Test API connectivity")
	fmt.Println("  interactive      - Start interactive test session")
	fmt.Println("  help             - Show this help message")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  XPAY_API_KEY     - Your X-Pay API key (required)")
	fmt.Println("  XPAY_MERCHANT_ID - Your merchant ID")
	fmt.Println("  XPAY_ENVIRONMENT - sandbox or live (default: sandbox)")
	fmt.Println("  XPAY_BASE_URL    - API base URL (default: auto-detected)")
	fmt.Println("  SERVER_PORT      - Server port (default: 3000)")
	fmt.Println("  XPAY_DEBUG       - Enable debug logging (default: false)")
}