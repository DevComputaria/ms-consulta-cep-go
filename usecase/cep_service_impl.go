package usecase

import (
	"example.com/hello/domain"
	"example.com/hello/interfaces/services"
)

// cepServiceImpl implements the CepService interface.
type cepServiceImpl struct {
	client services.ViaCepClient
}

// NewCepService creates a new instance of CepService.
// It takes a services.ViaCepClient as a dependency.
func NewCepService(client services.ViaCepClient) CepService {
	return &cepServiceImpl{
		client: client,
	}
}

// GetAddressByCep retrieves address details for a given CEP.
// It calls the FetchAddressFromViaCep method of the underlying ViaCepClient.
func (s *cepServiceImpl) GetAddressByCep(cep string) (*domain.Address, error) {
	// Basic validation could be added here if desired (e.g., length, numeric characters).
	// For now, we rely on the ViaCEP API for validation.

	address, err := s.client.FetchAddressFromViaCep(cep)
	if err != nil {
		return nil, err // Propagate the error from the client
	}
	return address, nil
}
