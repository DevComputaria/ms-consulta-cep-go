package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"example.com/hello/domain"
	"example.com/hello/usecase" // Will use usecase.CepServiceMock
)

func TestCepHandler_GetAddressByCepHandler(t *testing.T) {
	sampleAddress := &domain.Address{
		CEP:        "01001-000",
		Logradouro: "Praça da Sé",
		Bairro:     "Sé",
		Localidade: "São Paulo",
		UF:         "SP",
	}

	tests := []struct {
		name               string
		cepPath            string // The full path, e.g., "/cep/01001000" or "/cep/"
		mockAddress        *domain.Address
		mockServiceError   error
		expectedStatusCode int
		expectedBody       interface{} // Can be domain.Address or map[string]string for errors
		expectedHeaders    map[string]string
	}{
		{
			name:               "Successful Response",
			cepPath:            "/cep/01001000",
			mockAddress:        sampleAddress,
			mockServiceError:   nil,
			expectedStatusCode: http.StatusOK,
			expectedBody:       sampleAddress,
			expectedHeaders:    map[string]string{"Content-Type": "application/json"},
		},
		{
			name:               "CEP Not Found - service returns 'not found' error",
			cepPath:            "/cep/99999999",
			mockAddress:        nil,
			mockServiceError:   fmt.Errorf("address not found for CEP: 99999999"), // Error containing "not found"
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       map[string]string{"error": "Address not found for CEP: 99999999"},
			expectedHeaders:    map[string]string{"Content-Type": "application/json"},
		},
		{
			name:               "CEP Not Found - service returns 'failed to decode' error (treated as not found)",
			cepPath:            "/cep/88888888",
			mockAddress:        nil,
			mockServiceError:   fmt.Errorf("failed to decode response body"), // Error containing "failed to decode"
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       map[string]string{"error": "Address not found for CEP: 88888888"},
			expectedHeaders:    map[string]string{"Content-Type": "application/json"},
		},
		{
			name:               "Internal Server Error - generic service error",
			cepPath:            "/cep/12345678",
			mockAddress:        nil,
			mockServiceError:   errors.New("some internal service error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Internal server error"},
			expectedHeaders:    map[string]string{"Content-Type": "application/json"},
		},
		{
			name:               "Invalid CEP in path - empty CEP",
			cepPath:            "/cep/", // Empty CEP
			mockAddress:        nil,     // Service not called
			mockServiceError:   nil,     // Service not called
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "CEP must be provided in the URL path, e.g., /cep/01001000"},
			// Content-Type might not be application/json for http.Error default, let's check if handler sets it
		},
		{
			name:               "Invalid CEP in path - no CEP segment",
			cepPath:            "/cep", // No trailing slash, no CEP
			mockAddress:        nil,    // Service not called
			mockServiceError:   nil,    // Service not called
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "CEP must be provided in the URL path, e.g., /cep/01001000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Create mock service
			mockService := &usecase.CepServiceMock{
				MockAddress: tt.mockAddress,
				MockError:   tt.mockServiceError,
			}

			// Setup: Create handler with mock service
			handler := NewCepHandler(mockService)

			// Setup: Create request and response recorder
			req, err := http.NewRequest("GET", tt.cepPath, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			rr := httptest.NewRecorder()

			// Execute: Call the handler function
			handler.GetAddressByCepHandler(rr, req)

			// Assert: Check status code
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatusCode)
				t.Logf("Response body: %s", rr.Body.String())
			}

			// Assert: Check headers
			for key, expectedValue := range tt.expectedHeaders {
				if value := rr.Header().Get(key); value != expectedValue {
					t.Errorf("handler returned wrong header %s: got %q want %q", key, value, expectedValue)
				}
			}

			// Assert: Check response body
			// For JSON bodies, unmarshal and compare. For plain text (like http.Error might produce by default), direct string compare.
			if tt.expectedHeaders["Content-Type"] == "application/json" {
				var actualBody interface{}
				// Determine the type of expectedBody to unmarshal into the correct struct
				if _, ok := tt.expectedBody.(*domain.Address); ok {
					actualBody = &domain.Address{}
				} else if _, ok := tt.expectedBody.(map[string]string); ok {
					actualBody = make(map[string]string)
				} else {
					t.Fatalf("Unsupported type for expectedBody: %T for test %s", tt.expectedBody, tt.name)
				}
				
				err := json.Unmarshal(rr.Body.Bytes(), actualBody)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v. Body: %s", err, rr.Body.String())
				}

				if !reflect.DeepEqual(actualBody, tt.expectedBody) {
					// Try to provide more specific diff for maps
					if expectedMap, okE := tt.expectedBody.(map[string]string); okE {
						if actualMap, okA := actualBody.(map[string]string); okA {
							for k, vE := range expectedMap {
								if vA, ok := actualMap[k]; !ok || vA != vE {
									t.Errorf("handler returned unexpected body for key %s: got %v want %v", k, vA, vE)
								}
							}
						} else {
							t.Errorf("handler returned unexpected body type: got %T want map[string]string", actualBody)
						}
					} else {
						t.Errorf("handler returned unexpected body: got %v want %v", actualBody, tt.expectedBody)
					}
				}

			} else if tt.expectedBody != nil { // For non-JSON or if specific string error is expected
				expectedBodyStr, ok := tt.expectedBody.(string)
				if !ok {
					// If it's a map[string]string error, encode it to JSON string for comparison
					// This handles the case where http.Error is used and Content-Type isn't explicitly application/json
					// but the handler *does* write a JSON string.
					if errorMap, isMap := tt.expectedBody.(map[string]string); isMap {
						expectedBytes, _ := json.Marshal(errorMap)
						expectedBodyStr = string(expectedBytes)
						// The actual body from http.Error might have a newline
						if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expectedBodyStr) {
							t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expectedBodyStr)
						}
					} else {
						t.Fatalf("expectedBody is not a string and Content-Type is not application/json for test %s", tt.name)
					}
				} else {
					if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expectedBodyStr) {
						t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expectedBodyStr)
					}
				}
			}
		})
	}
}
