package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"example.com/hello/domain"
)

func TestViaCepClientImpl_FetchAddressFromViaCep(t *testing.T) {
	sampleAddress := &domain.Address{
		CEP:        "01001-000",
		Logradouro: "Praça da Sé",
		Complemento: "lado ímpar",
		Bairro:     "Sé",
		Localidade: "São Paulo",
		UF:         "SP",
		IBGE:       "3550308",
		GIA:        "1004",
		DDD:        "11",
		SIAFI:      "7107",
	}
	sampleAddressJSON, _ := json.Marshal(sampleAddress)

	tests := []struct {
		name           string
		cep            string
		serverHandler  func(w http.ResponseWriter, r *http.Request)
		expectedAddr   *domain.Address
		expectError    bool
		errorContains  string // Substring to check for in the error message
	}{
		{
			name: "Successful API Response",
			cep:  "01001000",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.URL.Path, "/01001000/json/") {
					http.Error(w, "Unexpected CEP in request URL", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(sampleAddressJSON)
			},
			expectedAddr:   sampleAddress,
			expectError:    false,
		},
		{
			name: "API Returns 404 Not Found",
			cep:  "99999999",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "CEP not found"}`))
			},
			expectedAddr:   nil,
			expectError:    true,
			errorContains:  "request failed with status code: 404",
		},
		{
			name: "API Returns Malformed JSON",
			cep:  "12345000",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"cep": "12345-000", "logradouro": "Rua Teste",`)) // Malformed
			},
			expectedAddr:   nil,
			expectError:    true,
			errorContains:  "failed to decode response body",
		},
		{
			name: "ViaCEP Not Found Response (empty CEP field)",
			cep:  "00000000",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				// ViaCEP returns 200 OK but with an empty/null CEP or an "erro: true" field
				// Based on current implementation, we check for empty CEP string in response
				w.Write([]byte(`{"cep": "", "logradouro": "", "uf": ""}`))
			},
			expectedAddr:   nil,
			expectError:    true,
			errorContains:  "address not found for CEP: 00000000",
		},
		{
			name: "ViaCEP Not Found Response (erro: true field - though current code does not check this explicitly)",
			cep:  "11111111",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				// Some ViaCEP responses might include {"erro": true}
				// Current implementation relies on empty address.CEP
				w.Write([]byte(`{"erro": true, "cep": ""}`))
			},
			expectedAddr:   nil,
			expectError:    true,
			errorContains:  "address not found for CEP: 11111111",
		},
		{
			name: "HTTP request creation failure (simulated by providing bad URL in client code - not directly testable here without altering tested code)",
			// This case is hard to test directly without injecting an error into http.NewRequest
			// or providing a malformed URL pattern to the client, which is not what we are testing.
			// The existing error handling for this in the SUT is `fmt.Errorf("failed to create request: %w", err)`
			// We will skip directly testing this specific internal error path.
			cep: "bad-request", // This won't cause NewRequest to fail, but illustrates intent.
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				// This handler should ideally not be called if NewRequest fails.
				http.Error(w, "Should not be reached", http.StatusInternalServerError)
			},
			// If we could force http.NewRequest to fail, we'd expect an error like:
			// expectError:    true,
			// errorContains:  "failed to create request",
			// For now, this test will behave like a normal "not found" or whatever the API returns for "bad-request"
			// Depending on ViaCEP, "bad-request" might be a 400 or other error.
			// Let's assume it's a 400 for this hypothetical scenario.
			expectedAddr:   nil,
			expectError:    true,
			errorContains:  "request failed with status code: 400", // Assuming ViaCEP returns 400 for completely invalid CEP format
		},
		{
			name: "HTTP client execution failure (e.g. network error)",
			// This is tested by shutting down the mock server before the client makes a request.
			cep:  "12312312",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				// Handler will be set up, but server shut down.
			},
			expectedAddr:   nil,
			expectError:    true,
			errorContains:  "failed to execute request:", // Error will contain more details from net/http
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverHandler))
			
			// Special case for testing http client execution failure
			if tt.name == "HTTP client execution failure (e.g. network error)" || tt.name == "HTTP request creation failure (simulated by providing bad URL in client code - not directly testable here without altering tested code)" {
				// For "failed to execute request", close the server immediately.
				// For the "bad-request" (which we're treating as a 400), we modify the handler.
				if tt.name == "HTTP client execution failure (e.g. network error)" {
					server.Close()
				} else { // "bad-request" test case
					// Re-assign handler for this specific sub-test to return 400
					// This is a bit of a workaround because the table-driven test setup is static.
					// A more complex setup might initialize the server per test case inside the loop.
					// For now, we'll just make sure this CEP actually results in a 400.
					// The actual ViaCEP API might return 400 for "bad-request".
					// Our mock server will simulate this.
					customServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest) // 400
					}))
					defer customServer.Close()
					server = customServer // use this server for this specific test
				}
			} else {
				defer server.Close()
			}


			// Use the NewViaCepClientWithHttpClient constructor with the test server's client
			client := NewViaCepClientWithHttpClient(server.Client())

			addr, err := client.FetchAddressFromViaCep(tt.cep)

			if tt.expectError {
				if err == nil {
					t.Errorf("FetchAddressFromViaCep() expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("FetchAddressFromViaCep() error = %q, expected to contain %q", err.Error(), tt.errorContains)
				}
			} else if err != nil {
				t.Errorf("FetchAddressFromViaCep() unexpected error: %v", err)
			}

			if !reflect.DeepEqual(addr, tt.expectedAddr) {
				t.Errorf("FetchAddressFromViaCep() address = %v, want %v", addr, tt.expectedAddr)
			}
		})
	}
}
