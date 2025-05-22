package main

import (
	"log"
	"net/http"

	httpHandler "example.com/hello/interfaces/http" // Alias for clarity
	"example.com/hello/interfaces/services"
	"example.com/hello/usecase"
)

func main() {
	// 1. Initialize the ViaCepClient
	viaCepClient := services.NewViaCepClient()

	// 2. Initialize the CepService
	cepService := usecase.NewCepService(viaCepClient)

	// 3. Initialize the CepHandler
	cepHandler := httpHandler.NewCepHandler(cepService)

	// 4. Register the HTTP handler function
	// This will handle requests like /cep/01001000, /cep/90210000, etc.
	// The handler itself will parse the CEP from the path.
	http.HandleFunc("/cep/", cepHandler.GetAddressByCepHandler)

	// 5. Start the HTTP server
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
