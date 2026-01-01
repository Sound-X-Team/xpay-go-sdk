# Changelog

All notable changes to the X-Pay Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-20

### Added
- Initial release of X-Pay Go SDK
- Support for Stripe payments
- Support for Mobile Money (Liberia) payments
- Support for X-Pay Wallet payments
- Customer management (CRUD operations)
- Webhook management and signature verification
- Currency utilities with decimal precision
- Comprehensive error handling
- Environment auto-detection from API keys
- Full context.Context support
- Type-safe API with Go structs
- Complete test suite
- Working examples for all major features
- Makefile for easy development

### Features
- **Payment Processing**: Create, retrieve, list, and cancel payments
- **Customer Management**: Full CRUD operations for customers
- **Webhook Management**: Create, update, delete, and test webhook endpoints
- **Currency Handling**: Proper decimal precision with shopspring/decimal
- **Error Handling**: Structured error responses with codes and details
- **Type Safety**: Full Go struct definitions with JSON tags
- **Testing**: Comprehensive unit tests and integration test suite

### Supported Payment Methods
- Stripe (USD, EUR, GBP, GHS)
- Mobile Money Liberia (USD)
- X-Pay Wallet (USD, GHS, EUR)

### Examples
- Basic payment creation with different payment methods
- Customer management operations
- Webhook setup and signature verification
- Error handling patterns
- Currency conversion utilities

### Development Tools
- Makefile with common development tasks
- Comprehensive test suite
- Code formatting and linting support
- Multi-platform build support