package papaya

import "encoding/json"

// CheckoutItem representa un producto o servicio a cobrar.
type CheckoutItem struct {
	Name               string `json:"name"`
	Quantity           int    `json:"quantity"`
	UnitPriceVESCents  int64  `json:"unit_price_ves_cents"`
	UnitPriceUSDCCents int64  `json:"unit_price_usdc_cents"`
}

// CheckoutItemsList maneja el caso donde la API devuelve los items
// como un JSON Array directo o como un string escapado.
type CheckoutItemsList []CheckoutItem

// UnmarshalJSON implementa el unmarshaling customizado para soportar strings escapados
func (l *CheckoutItemsList) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		// La API devolvió un string escapado, así que lo parseamos internamente
		return json.Unmarshal([]byte(str), (*[]CheckoutItem)(l))
	}
	
	// Si no era un string, lo parseamos como un array normal.
	// Usamos un tipo alias para evitar recursión infinita en este método.
	type Alias CheckoutItemsList
	return json.Unmarshal(data, (*Alias)(l))
}

// Checkout representa una intención de pago generada en Papaya.
type Checkout struct {
	ID               *string            `json:"id,omitempty"`
	AppID            *string            `json:"app_id,omitempty"`
	ExternalRef      *string            `json:"external_reference,omitempty"`
	AmountVESCents   *int64             `json:"amount_ves_cents,omitempty"`
	AmountUSDCCents  *int64             `json:"amount_usdc_cents,omitempty"`
	PaidCurrency     *string            `json:"paid_currency,omitempty"`
	Status           *string            `json:"status,omitempty"`
	Items            *CheckoutItemsList `json:"items,omitempty"`
	ExpiresAt        *int64             `json:"expires_at,omitempty"`
	CreatedAt        *int64             `json:"created_at,omitempty"`
}

// CreateCheckoutReq representa el payload para crear un checkout.
type CreateCheckoutReq struct {
	ExternalRef       string            `json:"external_reference,omitempty"`
	AmountVESCents    *int64            `json:"amount_ves_cents,omitempty"`
	AmountUSDCCents   *int64            `json:"amount_usdc_cents,omitempty"`
	Items             []CheckoutItem    `json:"items"`
	ExpirationMinutes *int              `json:"expiration_minutes,omitempty"`
}

// ListCheckoutsParams parámetros opcionales para paginar checkouts.
type ListCheckoutsParams struct {
	Limit           *int
	CursorCreatedAt *int64
	CursorID        *string
}
