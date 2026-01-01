package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/test-app/config"
	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
	"github.com/gorilla/mux"
)

// Server holds the HTTP server configuration
type Server struct {
	client *xpay.Client
	config *config.Config
	router *mux.Router
}

// Run starts the HTTP server 
func Run(cfg *config.Config) error {
	client := xpay.NewClient(cfg.XPay)
	
	server := &Server{
		client: client,
		config: cfg,
		router: mux.NewRouter(),
	}

	server.setupRoutes()

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	
	if cfg.Debug {
		cfg.Print()
	}

	fmt.Printf("🌐 Starting HTTP server on http://localhost%s\n", addr)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /health                    - Health check")
	fmt.Println("  GET  /health/xpay               - X-Pay API health check")
	fmt.Println("  POST /api/payments              - Create payment")
	fmt.Println("  GET  /api/payments              - List payments")
	fmt.Println("  GET  /api/payments/{id}         - Get payment")
	fmt.Println("  POST /api/payments/{id}/cancel  - Cancel payment")
	fmt.Println("  POST /api/customers             - Create customer")
	fmt.Println("  GET  /api/customers             - List customers")
	fmt.Println("  GET  /api/customers/{id}        - Get customer")
	fmt.Println("  PUT  /api/customers/{id}        - Update customer")
	fmt.Println("  DELETE /api/customers/{id}      - Delete customer")
	fmt.Println("  POST /api/webhooks              - Create webhook")
	fmt.Println("  GET  /api/webhooks              - List webhooks")
	fmt.Println("  GET  /api/webhooks/{id}         - Get webhook")
	fmt.Println("  PUT  /api/webhooks/{id}         - Update webhook")
	fmt.Println("  DELETE /api/webhooks/{id}       - Delete webhook")
	fmt.Println()

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      server.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return httpServer.ListenAndServe()
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)

	// Health endpoints
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
	s.router.HandleFunc("/health/xpay", s.handleXPayHealth).Methods("GET")

	// API routes
	api := s.router.PathPrefix("/api").Subrouter()
	
	// Payment routes
	api.HandleFunc("/payments", s.handleCreatePayment).Methods("POST")
	api.HandleFunc("/payments", s.handleListPayments).Methods("GET")
	api.HandleFunc("/payments/{id}", s.handleGetPayment).Methods("GET")
	api.HandleFunc("/payments/{id}/cancel", s.handleCancelPayment).Methods("POST")

	// Customer routes
	api.HandleFunc("/customers", s.handleCreateCustomer).Methods("POST")
	api.HandleFunc("/customers", s.handleListCustomers).Methods("GET")
	api.HandleFunc("/customers/{id}", s.handleGetCustomer).Methods("GET")
	api.HandleFunc("/customers/{id}", s.handleUpdateCustomer).Methods("PUT")
	api.HandleFunc("/customers/{id}", s.handleDeleteCustomer).Methods("DELETE")

	// Webhook routes
	api.HandleFunc("/webhooks", s.handleCreateWebhook).Methods("POST")
	api.HandleFunc("/webhooks", s.handleListWebhooks).Methods("GET")
	api.HandleFunc("/webhooks/{id}", s.handleGetWebhook).Methods("GET")
	api.HandleFunc("/webhooks/{id}", s.handleUpdateWebhook).Methods("PUT")
	api.HandleFunc("/webhooks/{id}", s.handleDeleteWebhook).Methods("DELETE")

	// Root endpoint
	s.router.HandleFunc("/", s.handleRoot).Methods("GET")
}