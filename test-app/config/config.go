package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
)

// Config holds all configuration for the test application
type Config struct {
	// X-Pay SDK configuration
	XPay *xpay.Config

	// Server configuration
	ServerPort int
	Debug      bool

	// Test configuration
	TestTimeout time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Required: API Key
	apiKey := os.Getenv("XPAY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("XPAY_API_KEY environment variable is required")
	}

	// Optional: Merchant ID
	merchantID := os.Getenv("XPAY_MERCHANT_ID")
	if merchantID == "" {
		// Use default test merchant ID
		merchantID = "548d8033-fbe9-411b-991f-f159cdee7745"
	}

	// Optional: Environment
	environment := xpay.Environment(os.Getenv("XPAY_ENVIRONMENT"))
	if environment == "" {
		environment = xpay.EnvironmentSandbox
	}

	// Optional: Base URL
	baseURL := os.Getenv("XPAY_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000" // Default for local development
	}

	// Optional: Server Port
	serverPort := 3000
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			serverPort = port
		}
	}

	// Optional: Debug mode
	debug := false
	if debugStr := os.Getenv("XPAY_DEBUG"); debugStr != "" {
		debug, _ = strconv.ParseBool(debugStr)
	}

	// Optional: Test timeout
	testTimeout := 30 * time.Second
	if timeoutStr := os.Getenv("TEST_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			testTimeout = timeout
		}
	}

	return &Config{
		XPay: &xpay.Config{
			APIKey:      apiKey,
			MerchantID:  merchantID,
			Environment: environment,
			BaseURL:     baseURL,
			Timeout:     30 * time.Second,
		},
		ServerPort:  serverPort,
		Debug:       debug,
		TestTimeout: testTimeout,
	}, nil
}

// Print prints the current configuration (without sensitive data)
func (c *Config) Print() {
	fmt.Println("Configuration:")
	fmt.Printf("  API Key: %s...%s\n", c.XPay.APIKey[:10], c.XPay.APIKey[len(c.XPay.APIKey)-4:])
	fmt.Printf("  Merchant ID: %s\n", c.XPay.MerchantID)
	fmt.Printf("  Environment: %s\n", c.XPay.Environment)
	fmt.Printf("  Base URL: %s\n", c.XPay.BaseURL)
	fmt.Printf("  Server Port: %d\n", c.ServerPort)
	fmt.Printf("  Debug Mode: %v\n", c.Debug)
	fmt.Printf("  Test Timeout: %v\n", c.TestTimeout)
	fmt.Println()
}