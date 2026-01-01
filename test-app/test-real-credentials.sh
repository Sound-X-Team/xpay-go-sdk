#!/bin/bash

# X-Pay Go SDK Test Application - Real Credentials Test Script

set -e

echo "🔐 X-Pay Go SDK - Real Credentials Test"
echo "======================================="

# Check if real API key is provided
if [ -z "$XPAY_REAL_API_KEY" ]; then
    echo "❌ XPAY_REAL_API_KEY environment variable is required"
    echo "   Set it with your real X-Pay API key:"
    echo "   export XPAY_REAL_API_KEY='sk_live_your_real_key_here'"
    echo "   export XPAY_REAL_MERCHANT_ID='your_real_merchant_id'"
    echo "   ./test-real-credentials.sh"
    exit 1
fi

# Check if real merchant ID is provided
if [ -z "$XPAY_REAL_MERCHANT_ID" ]; then
    echo "❌ XPAY_REAL_MERCHANT_ID environment variable is required"
    echo "   Set it with your real merchant ID:"
    echo "   export XPAY_REAL_MERCHANT_ID='your_real_merchant_id'"
    exit 1
fi

# Determine environment from API key
if [[ $XPAY_REAL_API_KEY == sk_live_* ]]; then
    ENVIRONMENT="live"
    BASE_URL="https://api.xpay-bits.com"
else
    ENVIRONMENT="sandbox"
    BASE_URL="https://api-sandbox.xpay-bits.com"
fi

# Override environment variables
export XPAY_API_KEY="$XPAY_REAL_API_KEY"
export XPAY_MERCHANT_ID="$XPAY_REAL_MERCHANT_ID"
export XPAY_ENVIRONMENT="$ENVIRONMENT"
export XPAY_BASE_URL="$BASE_URL"
export XPAY_DEBUG="true"

echo "🔧 Configuration:"
echo "   API Key: ${XPAY_API_KEY:0:10}...${XPAY_API_KEY: -4}"
echo "   Merchant ID: $XPAY_MERCHANT_ID"
echo "   Environment: $XPAY_ENVIRONMENT"
echo "   Base URL: $XPAY_BASE_URL"
echo

# Confirm before running with real credentials
echo "⚠️  WARNING: This will run tests with REAL API credentials!"
echo "   This may create actual payments, customers, and webhooks."
echo "   Make sure you understand the implications."
echo

read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Test cancelled"
    exit 1
fi

# Install dependencies and build
echo "📦 Installing dependencies..."
go mod download
go mod tidy

echo "🔨 Building application..."
go build -o bin/test-app .

# Run the test suite
echo "🧪 Running test suite with real credentials..."
echo "=============================================="

# First test connectivity
echo "📡 Testing API connectivity..."
if ./bin/test-app ping; then
    echo "✅ API connectivity successful"
else
    echo "❌ API connectivity failed"
    exit 1
fi

echo

# Run specific tests based on argument
case "${1:-basic}" in
    "basic")
        echo "🧪 Running basic test suite..."
        ./bin/test-app test
        ;;
    "payments")
        echo "💳 Running payment tests..."
        ./bin/test-app test-payments
        ;;
    "customers")
        echo "👥 Running customer tests..."
        ./bin/test-app test-customers
        ;;
    "webhooks")
        echo "🔗 Running webhook tests..."
        ./bin/test-app test-webhooks
        ;;
    "server")
        echo "🌐 Starting HTTP server with real credentials..."
        echo "   Server will be available at http://localhost:3000"
        echo "   Press Ctrl+C to stop"
        ./bin/test-app server
        ;;
    "interactive")
        echo "🎮 Starting interactive mode with real credentials..."
        ./bin/test-app interactive
        ;;
    *)
        echo "❌ Unknown test type: $1"
        echo "Available options: basic, payments, customers, webhooks, server, interactive"
        exit 1
        ;;
esac

echo
echo "🎉 Real credentials test completed!"
echo "   Remember to check your X-Pay dashboard for any created resources."