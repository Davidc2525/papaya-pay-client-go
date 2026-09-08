package papaya

import "errors"

var (
	// ErrInvalidSignatureFormat indica que el header X-Papaya-Signature no tiene el formato esperado (t=...,v1=...)
	ErrInvalidSignatureFormat = errors.New("formato de firma inválido: se esperaba t=<timestamp>,v1=<firma>")

	// ErrInvalidTimestamp indica que el timestamp proveído en la firma no es un número válido
	ErrInvalidTimestamp = errors.New("el timestamp de la firma no es válido")

	// ErrToleranceExceeded indica que la diferencia entre el tiempo actual y el de la firma supera la tolerancia permitida
	ErrToleranceExceeded = errors.New("el timestamp de la firma excede la tolerancia permitida")

	// ErrSignatureMismatch indica que la firma calculada no coincide con la proveída (posible alteración o clave incorrecta)
	ErrSignatureMismatch = errors.New("la firma proveída no coincide con el payload")

	// ErrInvalidPayload indica que el cuerpo del webhook no es un JSON válido o no coincide con la estructura esperada
	ErrInvalidPayload = errors.New("el cuerpo del webhook no es un JSON válido")
)
