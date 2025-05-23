package usecase

import (
	"example.com/hello/domain"
)

// CepServiceMock is a mock implementation of the CepService interface.
// It is used for testing purposes, particularly for the HTTP handler tests.
type CepServiceMock struct {
	MockAddress *domain.Address
	MockError   error
}

// GetAddressByCep mocks the behavior of fetching an address by CEP.
// It returns the pre-configured MockAddress and MockError.
func (m *CepServiceMock) GetAddressByCep(cep string) (*domain.Address, error) {
	return m.MockAddress, m.MockError
}

// NewCepServiceMock creates a new instance of CepServiceMock.
// This helper function can be used to easily set up the mock.
func NewCepServiceMock(address *domain.Address, err error) *CepServiceMock {
	return &CepServiceMock{
		MockAddress: address,
		MockError:   err,
	}
}
