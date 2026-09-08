package papaya

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client representa el cliente HTTP para interactuar con la API de Papaya
type Client struct {
	apiKey     string
	baseURL    string
	HTTPClient *http.Client
	
	// Servicios disponibles
	Checkouts *CheckoutsService
}

// NewClient crea un nuevo cliente API de Papaya.
func NewClient(apiKey string) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    "http://localhost:80",
		HTTPClient: http.DefaultClient,
	}
	c.Checkouts = &CheckoutsService{client: c}
	return c
}

// SetBaseURL permite sobreescribir la URL de la API (ej. para Sandbox)
func (c *Client) SetBaseURL(url string) {
	c.baseURL = strings.TrimSuffix(url, "/")
}

// doRequest es un método interno para ejecutar peticiones HTTP.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, responseObj interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Papaya API Error (status %d): %s", resp.StatusCode, string(errBody))
	}

	if responseObj != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseObj); err != nil {
			return err
		}
	}

	return nil
}
