package papaya

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DefaultTolerance es la tolerancia por defecto para validar timestamps (5 minutos).
const DefaultTolerance = 5 * time.Minute

// Transaction representa la información monetaria del pago.
type Transaction struct {
	ID             string `json:"id"`
	FeeCents       int64  `json:"fee_cents"`
	NetAmountCents int64  `json:"net_amount_cents"`
	PaidCurrency   string `json:"paid_currency"`
}

// WebhookEvent representa la estructura del JSON enviado por Papaya.
type WebhookEvent struct {
	AppID      string      `json:"app_id"`
	CheckoutID string      `json:"checkout_id"`
	Status     string      `json:"status"`
	Transaction Transaction `json:"transaction"`
}

// ConstructEvent valida la firma criptográfica y parsea el payload crudo en un WebhookEvent.
//
// Parámetros:
//   - rawPayload: El cuerpo completo y exacto de la petición HTTP recibida.
//   - signatureHeader: El valor del header HTTP 'X-Papaya-Signature'.
//   - secret: El webhook secret del comercio.
//   - tolerance: Tiempo máximo permitido de diferencia entre el timestamp de la firma y el actual.
func ConstructEvent(rawPayload []byte, signatureHeader string, secret string, tolerance time.Duration) (*WebhookEvent, error) {
	// 1. Extraer timestamp (t) y firma (v1)
	parts := strings.Split(signatureHeader, ",")
	var tStr, v1Str string

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			switch kv[0] {
			case "t":
				tStr = kv[1]
			case "v1":
				v1Str = kv[1]
			}
		}
	}

	if tStr == "" || v1Str == "" {
		return nil, ErrInvalidSignatureFormat
	}

	timestamp, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidTimestamp
	}

	// 2. Validar Tolerancia
	now := time.Now().Unix()
	diff := now - timestamp
	if diff < 0 {
		diff = -diff
	}

	if time.Duration(diff)*time.Second > tolerance {
		return nil, ErrToleranceExceeded
	}

	// 3. Calcular el hash esperado: hmac_sha256(secret, t + "." + rawPayload)
	signedPayload := fmt.Sprintf("%d.%s", timestamp, rawPayload)
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// 4. Comparación Segura (Constant-Time Compare) para evitar Timing Attacks
	expectedMAC, _ := hex.DecodeString(expectedSignature)
	providedMAC, err := hex.DecodeString(v1Str)
	if err != nil {
		return nil, ErrSignatureMismatch
	}

	if !hmac.Equal(expectedMAC, providedMAC) {
		return nil, ErrSignatureMismatch
	}

	// 5. Parseo del JSON
	var event WebhookEvent
	if err := json.Unmarshal(rawPayload, &event); err != nil {
		return nil, ErrInvalidPayload
	}

	return &event, nil
}
