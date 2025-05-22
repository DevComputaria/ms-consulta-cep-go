package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/hello/usecase"
	"fmt" // Added for error checking
)

// CepHandler handles HTTP requests related to CEP information.
type CepHandler struct {
	service usecase.CepService
}

// NewCepHandler creates a new instance of CepHandler.
func NewCepHandler(service usecase.CepService) *CepHandler {
	return &CepHandler{
		service: service,
	}
}

// GetAddressByCepHandler handles the request to get an address by CEP.
// It expects the CEP to be part of the URL path, e.g., /cep/01001000.
func (h *CepHandler) GetAddressByCepHandler(w http.ResponseWriter, r *http.Request) {
	// Extract CEP from path, assuming path is /cep/{cepValue}
	// For a production system, a router like gorilla/mux would be better.
	cep := strings.TrimPrefix(r.URL.Path, "/cep/")
	if cep == "" || cep == r.URL.Path { // Check if TrimPrefix did anything
		http.Error(w, `{"error": "CEP must be provided in the URL path, e.g., /cep/01001000"}`, http.StatusBadRequest)
		return
	}

	address, err := h.service.GetAddressByCep(cep)
	if err != nil {
		// Check if the error message indicates "not found"
		// This is a simple check. In a real application, custom error types or codes would be better.
		if strings.Contains(strings.ToLower(err.Error()), "not found") || strings.Contains(strings.ToLower(err.Error()), "failed to decode response body") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Address not found for CEP: %s", cep)})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(address); err != nil {
		// If encoding fails, it's an internal server error, though the headers might already be sent.
		// Log this error for server-side diagnostics.
		// For the client, it might be too late to send a different status code.
		fmt.Printf("Error encoding address to JSON: %v\n", err) // Log to server console
	}
}
