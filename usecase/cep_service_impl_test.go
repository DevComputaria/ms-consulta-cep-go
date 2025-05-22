package usecase

import (
	"errors"
	"testing"
	"reflect" // For deep equality comparison

	"example.com/hello/domain"
	"example.com/hello/interfaces/services"
)

func TestCepServiceImpl_GetAddressByCep(t *testing.T) {
	sampleAddress := &domain.Address{
		CEP:        "01001-000",
		Logradouro: "Praça da Sé",
		Bairro:     "Sé",
		Localidade: "São Paulo",
		UF:         "SP",
	}
	sampleError := errors.New("client error")

	tests := []struct {
		name          string
		cep           string
		mockAddress   *domain.Address
		mockError     error
		expectedAddr  *domain.Address
		expectedError error
	}{
		{
			name:          "Successful fetch",
			cep:           "01001000",
			mockAddress:   sampleAddress,
			mockError:     nil,
			expectedAddr:  sampleAddress,
			expectedError: nil,
		},
		{
			name:          "Error from client",
			cep:           "12345678",
			mockAddress:   nil,
			mockError:     sampleError,
			expectedAddr:  nil,
			expectedError: sampleError,
		},
		{
			name:          "Client returns nil address and nil error (unexpected but testable)",
			cep:           "00000000",
			mockAddress:   nil,
			mockError:     nil,
			expectedAddr:  nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Create the mock client
			mockClient := &services.ViaCepClientMock{
				MockAddress: tt.mockAddress,
				MockError:   tt.mockError,
			}

			// Setup: Create the service implementation with the mock client
			service := NewCepService(mockClient)

			// Execute: Call the method being tested
			addr, err := service.GetAddressByCep(tt.cep)

			// Assert: Check the address
			if !reflect.DeepEqual(addr, tt.expectedAddr) {
				t.Errorf("GetAddressByCep() address = %v, want %v", addr, tt.expectedAddr)
			}

			// Assert: Check the error
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("GetAddressByCep() error = nil, want %v", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() { // Compare error messages
					t.Errorf("GetAddressByCep() error = %q, want %q", err.Error(), tt.expectedError.Error())
				}
			} else if err != nil {
				t.Errorf("GetAddressByCep() error = %v, want nil", err)
			}
		})
	}
}
