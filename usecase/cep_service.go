package usecase

import "example.com/hello/domain"

// CepService is an interface for fetching address information based on CEP.
type CepService interface {
	// GetAddressByCep retrieves address details for a given CEP.
	// It returns a pointer to an Address struct or an error if the CEP is not found or an issue occurs.
	GetAddressByCep(cep string) (*domain.Address, error)
}
