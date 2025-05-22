package services

import (
	"example.com/hello/domain"
)

// ViaCepClientMock is a mock implementation of the ViaCepClient interface.
// It is used for testing purposes.
type ViaCepClientMock struct {
	MockAddress *domain.Address
	MockError   error
}

// FetchAddressFromViaCep mocks the behavior of fetching an address from ViaCEP.
// It returns the pre-configured MockAddress and MockError.
func (m *ViaCepClientMock) FetchAddressFromViaCep(cep string) (*domain.Address, error) {
	return m.MockAddress, m.MockError
}

// NewViaCepClientMock creates a new instance of ViaCepClientMock.
// This is not strictly necessary as the struct can be initialized directly,
// but it can be useful for consistency or if more complex initialization is needed later.
func NewViaCepClientMock(address *domain.Address, err error) *ViaCepClientMock {
	return &ViaCepClientMock{
		MockAddress: address,
		MockError:   err,
	}
}
