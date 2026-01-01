# X-Pay Go SDK Test Application

A comprehensive test application that demonstrates the X-Pay Go SDK capabilities with a REST API server and CLI test runner.

## Features

- **REST API Server**: Complete HTTP server with payment, customer, and webhook endpoints
- **CLI Test Runner**: Command-line interface to test all SDK features
- **Health Checks**: Built-in health endpoints for monitoring
- **Error Handling**: Comprehensive error handling and responses
- **Real API Testing**: Scripts to test with real API credentials
- **Docker Support**: Containerized testing environment

## Structure

```
test-app/
├── README.md              # This file
├── go.mod                 # Go module for test app
├── main.go                # Main application entry point
├── server/                # HTTP server implementation
│   ├── server.go         # Server setup and routing
│   ├── handlers.go       # HTTP handlers
│   └── middleware.go     # HTTP middleware
├── cli/                   # CLI test runner
│   ├── cli.go            # CLI commands
│   └── tests.go          # Test suite runner
├── config/               # Configuration
│   └── config.go         # Configuration management
├── docker-compose.yml    # Docker composition
├── Dockerfile           # Docker container
├── run-tests.sh         # Test script
└── test-real-credentials.sh # Real API test script
```

## Quick Start

### Prerequisites

1. Go 1.19 or higher
2. X-Pay API credentials (sandbox or live)

### Installation

```bash
# Install dependencies
go mod download

# Build the application
go build -o bin/test-app
```

### Configuration

Set environment variables:

```bash
export XPAY_API_KEY="sk_sandbox_your_key_here"
export XPAY_MERCHANT_ID="your_merchant_id_here"
export XPAY_ENVIRONMENT="sandbox"
export XPAY_BASE_URL="http://localhost:8000"  # For local development
```

Or create a `.env` file:

```env
XPAY_API_KEY=sk_sandbox_your_key_here
XPAY_MERCHANT_ID=your_merchant_id_here
XPAY_ENVIRONMENT=sandbox
XPAY_BASE_URL=http://localhost:8000
```

## Running the Application

### 1. HTTP Server Mode

Start the REST API server:

```bash
go run main.go server
# or
./bin/test-app server
```

The server will start on `http://localhost:3000`

#### API Endpoints

**Payments**
- `POST /api/payments` - Create payment
- `GET /api/payments` - List payments
- `GET /api/payments/{id}` - Get payment
- `POST /api/payments/{id}/cancel` - Cancel payment

**Customers**
- `POST /api/customers` - Create customer
- `GET /api/customers` - List customers
- `GET /api/customers/{id}` - Get customer
- `PUT /api/customers/{id}` - Update customer
- `DELETE /api/customers/{id}` - Delete customer

**Webhooks**
- `POST /api/webhooks` - Create webhook
- `GET /api/webhooks` - List webhooks
- `GET /api/webhooks/{id}` - Get webhook
- `PUT /api/webhooks/{id}` - Update webhook
- `DELETE /api/webhooks/{id}` - Delete webhook

**Health**
- `GET /health` - Health check
- `GET /health/xpay` - X-Pay API connectivity check

### 2. CLI Test Mode

Run comprehensive SDK tests:

```bash
go run main.go test
# or
./bin/test-app test
```

Available CLI commands:
- `test` - Run full test suite
- `test-payments` - Test payment operations only
- `test-customers` - Test customer operations only
- `test-webhooks` - Test webhook operations only
- `ping` - Test API connectivity

### 3. Interactive Mode

Run interactive test session:

```bash
go run main.go interactive
```

## Example Usage

### Create a Payment via API

```bash
curl -X POST http://localhost:3000/api/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "29.99",
    "currency": "USD",
    "payment_method": "stripe",
    "description": "Test payment from Go SDK"
  }'
```

### Create a Customer via API

```bash
curl -X POST http://localhost:3000/api/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+231700000000"
  }'
```

### Check Health

```bash
curl http://localhost:3000/health
```

## Testing Scripts

### Run All Tests

```bash
./run-tests.sh
```

### Test with Real Credentials

```bash
./test-real-credentials.sh
```

## Docker Support

### Build and Run with Docker

```bash
# Build the image
docker build -t xpay-go-test-app .

# Run the container
docker run -p 3000:3000 -e XPAY_API_KEY=your_key_here xpay-go-test-app
```

### Using Docker Compose

```bash
docker-compose up --build
```

## Development

### Adding New Tests

Add test functions to `cli/tests.go`:

```go
func (t *TestRunner) TestNewFeature() error {
    // Your test implementation
    return nil
}
```

### Adding New API Endpoints

Add handlers to `server/handlers.go`:

```go
func (h *Handlers) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Your handler implementation
}
```

## Expected Output

When running tests, you should see output like:

```
🧪 X-Pay Go SDK Test Application
================================

✅ API Connectivity Test - PASSED
✅ Payment Creation Test - PASSED  
✅ Customer Management Test - PASSED
✅ Webhook Setup Test - PASSED
✅ Error Handling Test - PASSED

📊 Results: 5/5 tests passed (100%)
🎉 All tests completed successfully!
```

## Troubleshooting

### Common Issues

1. **API Key Not Set**: Ensure `XPAY_API_KEY` environment variable is set
2. **Server Not Running**: Make sure the X-Pay backend server is running on the configured base URL
3. **Network Issues**: Check firewall and network connectivity

### Debug Mode

Enable debug logging:

```bash
export XPAY_DEBUG=true
go run main.go test
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## Support

- Documentation: https://docs.xpay-bits.com
- Issues: GitHub Issues
- Support: support@xpay-bits.com