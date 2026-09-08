package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	papaya "github.com/Davidc2525/papaya-pay-client-go"
)

const AMOUNT_SCALE = 10000.0

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Obtener la firma del header
	signatureHeader := r.Header.Get("X-Papaya-Signature")
	if signatureHeader == "" {
		http.Error(w, "Falta el header X-Papaya-Signature", http.StatusBadRequest)
		return
	}

	// 2. Leer el cuerpo de la petición CRUDA en memoria.
	rawPayload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el cuerpo de la petición", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 3. Obtener el secreto y procesar el evento de webhook
	secret := os.Getenv("PAPAYA_WEBHOOK_SECRET")
	if secret == "" {
		secret = "mi_secreto_de_prueba"
	}

	event, err := papaya.ConstructEvent(rawPayload, signatureHeader, secret, papaya.DefaultTolerance)
	if err != nil {
		log.Printf("⚠️ Validación de Webhook fallida: %v\n", err)
		http.Error(w, fmt.Sprintf("Webhook Error: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Recibido evento seguro para checkout_id: %s", event.CheckoutID)
	log.Printf("Estado: %s", event.Status)

	if event.Status == "PAID" {
		// 4. Inicializar cliente API de Papaya
		apiKey := os.Getenv("PAPAYA_API_KEY")
		if apiKey == "" {
			apiKey = "sk_live_mi_api_key"
		}
		log.Printf("API Key: %s\n", apiKey)
		
		client := papaya.NewClient(apiKey)
		client.SetBaseURL("http://localhost:80") // Descomentar para desarrollo local
		
		// 5. Obtener los detalles completos del checkout
		log.Println("Consultando detalles del checkout a la API...")
		checkout, err := client.Checkouts.Get(r.Context(), event.CheckoutID)
		if err != nil {
			log.Printf("Error obteniendo detalles del checkout: %v\n", err)
			http.Error(w, "Error consultando la API", http.StatusInternalServerError)
			return
		}

		if checkout.ExternalRef != nil {
			log.Printf("Pedido Interno (external_reference): %s\n", *checkout.ExternalRef)
		}
		if checkout.Items != nil {
			log.Printf("Ítems pagados: %+v\n", *checkout.Items)
		}
		
		// TODO: Procesar el pago y entregar la orden
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"received": true}`))
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	port := ":8181"
	log.Printf("Servidor de ejemplo escuchando en el puerto %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
