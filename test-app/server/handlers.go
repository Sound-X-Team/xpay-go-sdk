package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Sound-X-Team/x-pay/integrations/sdks/golang/xpay"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// writeJSON writes a JSON response
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeSuccess writes a successful JSON response
func (s *Server) writeSuccess(w http.ResponseWriter, data interface{}) {
	s.writeJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// writeError writes an error JSON response
func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

// Root handler
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"name":        "X-Pay Go SDK Test Application",
		"version":     "1.0.0",
		"environment": s.config.XPay.Environment,
		"endpoints": map[string]string{
			"health":    "/health",
			"payments":  "/api/payments",
			"customers": "/api/customers",
			"webhooks":  "/api/webhooks",
		},
	}
	s.writeSuccess(w, info)
}

// Health check handler
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(time.Now()).String(),
	}
	s.writeSuccess(w, health)
}

// X-Pay API health check handler
func (s *Server) handleXPayHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ping, err := s.client.Ping(ctx)
	if err != nil {
		s.writeError(w, http.StatusServiceUnavailable, fmt.Sprintf("X-Pay API unavailable: %v", err))
		return
	}

	health := map[string]interface{}{
		"status":         "healthy",
		"xpay_api":       "connected",
		"xpay_success":   ping.Success,
		"xpay_timestamp": ping.Timestamp,
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	s.writeSuccess(w, health)
}

// Payment handlers

type CreatePaymentRequest struct {
	Amount            string                     `json:"amount"`
	Currency          string                     `json:"currency,omitempty"`
	PaymentMethod     string                     `json:"payment_method"`
	Description       string                     `json:"description,omitempty"`
	CustomerID        string                     `json:"customer_id,omitempty"`
	PaymentMethodData *xpay.PaymentMethodData    `json:"payment_method_data,omitempty"`
	Metadata          map[string]interface{}     `json:"metadata,omitempty"`
}

func (s *Server) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid amount format")
		return
	}

	paymentReq := &xpay.PaymentRequest{
		Amount:            amount,
		Currency:          xpay.Currency(req.Currency),
		PaymentMethod:     xpay.PaymentMethod(req.PaymentMethod),
		Description:       req.Description,
		CustomerID:        req.CustomerID,
		PaymentMethodData: req.PaymentMethodData,
		Metadata:          req.Metadata,
	}

	// Add server metadata
	if paymentReq.Metadata == nil {
		paymentReq.Metadata = make(map[string]interface{})
	}
	paymentReq.Metadata["created_via"] = "go-sdk-test-app-server"

	payment, err := s.client.Payments.Create(r.Context(), paymentReq)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, payment)
}

func (s *Server) handleListPayments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	req := &xpay.ListPaymentsRequest{}
	
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}
	
	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			req.Offset = o
		}
	}
	
	if status := query.Get("status"); status != "" {
		req.Status = xpay.PaymentStatus(status)
	}
	
	if customerID := query.Get("customer_id"); customerID != "" {
		req.CustomerID = customerID
	}

	payments, err := s.client.Payments.List(r.Context(), req)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, payments)
}

func (s *Server) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]

	payment, err := s.client.Payments.Retrieve(r.Context(), paymentID)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, payment)
}

func (s *Server) handleCancelPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]

	payment, err := s.client.Payments.Cancel(r.Context(), paymentID)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, payment)
}

// Customer handlers

func (s *Server) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req xpay.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Add server metadata
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["created_via"] = "go-sdk-test-app-server"

	customer, err := s.client.Customers.Create(r.Context(), &req)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, customer)
}

func (s *Server) handleListCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	req := &xpay.ListCustomersRequest{}
	
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}
	
	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			req.Offset = o
		}
	}
	
	if email := query.Get("email"); email != "" {
		req.Email = email
	}
	
	if phone := query.Get("phone"); phone != "" {
		req.Phone = phone
	}

	customers, err := s.client.Customers.List(r.Context(), req)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, customers)
}

func (s *Server) handleGetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	customer, err := s.client.Customers.Retrieve(r.Context(), customerID)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, customer)
}

func (s *Server) handleUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	var req xpay.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	customer, err := s.client.Customers.Update(r.Context(), customerID, &req)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, customer)
}

func (s *Server) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	err := s.client.Customers.Delete(r.Context(), customerID)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, map[string]string{"message": "Customer deleted successfully"})
}

// Webhook handlers

func (s *Server) handleCreateWebhook(w http.ResponseWriter, r *http.Request) {
	var req xpay.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	webhook, err := s.client.Webhooks.Create(r.Context(), &req)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, webhook)
}

func (s *Server) handleListWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := s.client.Webhooks.List(r.Context())
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, webhooks)
}

func (s *Server) handleGetWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webhookID := vars["id"]

	webhook, err := s.client.Webhooks.Retrieve(r.Context(), webhookID)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, webhook)
}

func (s *Server) handleUpdateWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webhookID := vars["id"]

	var req xpay.UpdateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	webhook, err := s.client.Webhooks.Update(r.Context(), webhookID, &req)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, webhook)
}

func (s *Server) handleDeleteWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webhookID := vars["id"]

	err := s.client.Webhooks.Delete(r.Context(), webhookID)
	if err != nil {
		if xpayErr, ok := err.(*xpay.Error); ok {
			s.writeError(w, xpayErr.StatusCode, xpayErr.Message)
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.writeSuccess(w, map[string]string{"message": "Webhook deleted successfully"})
}