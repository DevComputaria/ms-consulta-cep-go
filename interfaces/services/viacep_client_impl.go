package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/hello/domain"
)

// viaCepClientImpl implements the ViaCepClient interface.
type viaCepClientImpl struct {
	httpClient *http.Client
}

// NewViaCepClient creates a new instance of ViaCepClient with the default http client.
func NewViaCepClient() ViaCepClient {
	return &viaCepClientImpl{
		httpClient: http.DefaultClient,
	}
}

// NewViaCepClientWithHttpClient creates a new instance of ViaCepClient with a custom http client.
// This is useful for testing purposes.
func NewViaCepClientWithHttpClient(client *http.Client) ViaCepClient {
	return &viaCepClientImpl{
		httpClient: client,
	}
}

// FetchAddressFromViaCep fetches address details for a given CEP from the ViaCEP API.
func (c *viaCepClientImpl) FetchAddressFromViaCep(cep string) (*domain.Address, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var address domain.Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// ViaCEP returns a normal response with `cep: null` or `cep: ""` when the CEP is not found.
	// We check if the CEP field in the response is empty, which indicates that the address was not found.
	if address.CEP == "" {
		return nil, fmt.Errorf("address not found for CEP: %s", cep)
	}


	return &address, nil
}
