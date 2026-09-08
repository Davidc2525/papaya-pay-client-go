package papaya

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CheckoutsService agrupa los endpoints relacionados con Checkouts.
type CheckoutsService struct {
	client *Client
}

// Create genera una intención de pago (Checkout).
func (s *CheckoutsService) Create(ctx context.Context, req CreateCheckoutReq) (*Checkout, error) {
	var checkout Checkout
	err := s.client.doRequest(ctx, "POST", "/v1/gateway/api/checkouts", req, &checkout)
	return &checkout, err
}

// List obtiene el historial de checkouts del comercio.
func (s *CheckoutsService) List(ctx context.Context, params *ListCheckoutsParams) ([]Checkout, error) {
	path := "/v1/gateway/api/checkouts"
	if params != nil {
		query := url.Values{}
		if params.Limit != nil {
			query.Add("limit", strconv.Itoa(*params.Limit))
		}
		if params.CursorCreatedAt != nil {
			query.Add("cursor_created_at", strconv.FormatInt(*params.CursorCreatedAt, 10))
		}
		if params.CursorID != nil {
			query.Add("cursor_id", *params.CursorID)
		}
		encoded := query.Encode()
		if encoded != "" {
			path += "?" + encoded
		}
	}

	var checkouts []Checkout
	err := s.client.doRequest(ctx, "GET", path, nil, &checkouts)
	return checkouts, err
}

// Get obtiene los detalles completos de un checkout por su ID.
func (s *CheckoutsService) Get(ctx context.Context, id string) (*Checkout, error) {
	var checkout Checkout
	path := fmt.Sprintf("/v1/gateway/api/checkouts/%s", id)
	err := s.client.doRequest(ctx, "GET", path, nil, &checkout)
	return &checkout, err
}

// Cancel cancela un checkout pendiente.
func (s *CheckoutsService) Cancel(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v1/gateway/api/checkouts/%s/cancel", id)
	return s.client.doRequest(ctx, "POST", path, nil, nil)
}
