# Papaya Go SDK

SDK oficial de Papaya para Go (Golang). Incluye un cliente HTTP fuertemente tipado y funciones de seguridad criptográfica para la verificación de Webhooks.

## Uso de la Librería

### 1. Cliente API (Checkouts)
Importa el paquete y crea un cliente usando tu API Key. 

```go
package main

import (
	"context"
	"fmt"
	papaya "github.com/Davidc2525/papaya-pay-client-go" // Ajusta el import local o remoto
)

func main() {
	client := papaya.NewClient("sk_live_...")
	// client.SetBaseURL("http://localhost:80") // Opcional
	
	// Consultar un checkout
	checkout, err := client.Checkouts.Get(context.Background(), "ID_DEL_CHECKOUT")
	if err != nil {
		panic(err)
	}
	
	if checkout.Items != nil {
		fmt.Printf("Ítems comprados: %+v\n", *checkout.Items)
	}
}
```

### 2. Validación de Webhooks
```go
import papaya "github.com/papaya/papaya-pay-client"

// ... dentro de tu handler HTTP
event, err := papaya.ConstructEvent(
    rawBodyBytes, // El cuerpo crudo (io.ReadAll)
    signatureHeader, 
    "tu_webhook_secret", 
    papaya.DefaultTolerance,
)

if err != nil {
    // Manejar firma inválida
}
```

---

## Correr el Ejemplo (net/http)

El ejemplo incluido levanta un servidor HTTP usando la librería estándar de Go en el puerto `8181`. Te mostrará un flujo real donde recibe un webhook y luego llama al cliente API para conseguir más detalles.

1. Clona el repositorio y ve a la carpeta:
   ```bash
   git clone https://github.com/Davidc2525/papaya-pay-client-go.git
   cd papaya-pay-client-go
   ```
2. Asegúrate de tener las dependencias al día:
   ```bash
   go mod tidy
   ```
3. Ejecuta el servidor:
   ```bash
   go run examples/http/main.go
   ```
