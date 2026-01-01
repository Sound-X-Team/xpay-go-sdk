package xpay

import (
	"context"
	"fmt"
)

// CustomersService handles customer-related operations
type CustomersService struct {
	client *Client
}

// NewCustomersService creates a new customers service
func NewCustomersService(client *Client) *CustomersService {
	return &CustomersService{client: client}
}

// Create creates a new customer
func (s *CustomersService) Create(ctx context.Context, req *CreateCustomerRequest) (*Customer, error) {
	var result APIResponse[Customer]
	path := fmt.Sprintf("/v1/api/merchants/%s/customers", s.client.merchantID)
	
	err := s.client.httpClient.Post(ctx, path, req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Retrieve retrieves a customer by ID
func (s *CustomersService) Retrieve(ctx context.Context, customerID string) (*Customer, error) {
	var result APIResponse[Customer]
	path := fmt.Sprintf("/v1/api/merchants/%s/customers/%s", s.client.merchantID, customerID)
	
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Update updates a customer
func (s *CustomersService) Update(ctx context.Context, customerID string, req *UpdateCustomerRequest) (*Customer, error) {
	var result APIResponse[Customer]
	path := fmt.Sprintf("/v1/api/merchants/%s/customers/%s", s.client.merchantID, customerID)
	
	err := s.client.httpClient.Put(ctx, path, req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// List lists customers with optional filters
func (s *CustomersService) List(ctx context.Context, req *ListCustomersRequest) (*ListResponse[Customer], error) {
	params := make(map[string]interface{})
	
	if req != nil {
		if req.Limit > 0 {
			params["limit"] = req.Limit
		}
		if req.Offset > 0 {
			params["offset"] = req.Offset
		}
		if req.Email != "" {
			params["email"] = req.Email
		}
		if req.Phone != "" {
			params["phone"] = req.Phone
		}
	}

	queryString := BuildQueryString(params)
	path := fmt.Sprintf("/v1/api/merchants/%s/customers%s", s.client.merchantID, queryString)

	var result APIResponse[struct {
		Customers []Customer `json:"customers"`
		Total     int        `json:"total"`
	}]
	
	err := s.client.httpClient.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &ListResponse[Customer]{
		Items:  result.Data.Customers,
		Total:  result.Data.Total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// Delete deletes a customer
func (s *CustomersService) Delete(ctx context.Context, customerID string) error {
	path := fmt.Sprintf("/v1/api/merchants/%s/customers/%s", s.client.merchantID, customerID)
	
	var result APIResponse[interface{}]
	err := s.client.httpClient.Delete(ctx, path, &result)
	if err != nil {
		return err
	}

	return nil
}