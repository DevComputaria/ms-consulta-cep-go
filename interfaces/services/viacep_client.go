package services

import "example.com/hello/domain"

// ViaCepClient is an interface for interacting with the ViaCEP API.
type ViaCepClient interface {
	// FetchAddressFromViaCep fetches address details for a given CEP from the ViaCEP API.
	// It returns a pointer to an Address struct or an error if the CEP is not found or an issue occurs.
	FetchAddressFromViaCep(cep string) (*domain.Address, error)
}
