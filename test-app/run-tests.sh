#!/bin/bash

# X-Pay Go SDK Test Application - Test Runner Script

set -e

echo "🧪 X-Pay Go SDK Test Application"
echo "================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "📋 Go version: $GO_VERSION"

# Set default environment variables if not set
export XPAY_API_KEY="${XPAY_API_KEY:-sk_sandbox_7c845adf-f658-4f29-9857-7e8a8708}"
export XPAY_MERCHANT_ID="${XPAY_MERCHANT_ID:-548d8033-fbe9-411b-991f-f159cdee7745}"
export XPAY_ENVIRONMENT="${XPAY_ENVIRONMENT:-sandbox}"
export XPAY_BASE_URL="${XPAY_BASE_URL:-http://localhost:8000}"
export XPAY_DEBUG="${XPAY_DEBUG:-true}"

echo "🔧 Configuration:"
echo "   API Key: ${XPAY_API_KEY:0:10}...${XPAY_API_KEY: -4}"
echo "   Merchant ID: $XPAY_MERCHANT_ID"
echo "   Environment: $XPAY_ENVIRONMENT"
echo "   Base URL: $XPAY_BASE_URL"
echo "   Debug: $XPAY_DEBUG"
echo

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Build the application
echo "🔨 Building application..."
go build -o bin/test-app .

# Run tests based on argument
case "${1:-test}" in
    "server")
        echo "🌐 Starting HTTP server..."
        ./bin/test-app server
        ;;
    "test")
        echo "🧪 Running comprehensive test suite..."
        ./bin/test-app test
        ;;
    "test-payments")
        echo "💳 Running payment tests..."
        ./bin/test-app test-payments
        ;;
    "test-customers")
        echo "👥 Running customer tests..."
        ./bin/test-app test-customers
        ;;
    "test-webhooks")
        echo "🔗 Running webhook tests..."
        ./bin/test-app test-webhooks
        ;;
    "ping")
        echo "📡 Testing API connectivity..."
        ./bin/test-app ping
        ;;
    "interactive")
        echo "🎮 Starting interactive mode..."
        ./bin/test-app interactive
        ;;
    "clean")
        echo "🧹 Cleaning build artifacts..."
        rm -rf bin/
        go clean
        echo "✅ Clean completed"
        ;;
    "help")
        echo "Available commands:"
        echo "  server           - Start HTTP server"
        echo "  test             - Run comprehensive test suite"
        echo "  test-payments    - Run payment tests only"
        echo "  test-customers   - Run customer tests only"
        echo "  test-webhooks    - Run webhook tests only"
        echo "  ping             - Test API connectivity"
        echo "  interactive      - Start interactive mode"
        echo "  clean            - Clean build artifacts"
        echo "  help             - Show this help"
        ;;
    *)
        echo "❌ Unknown command: $1"
        echo "Run './run-tests.sh help' for available commands"
        exit 1
        ;;
esac