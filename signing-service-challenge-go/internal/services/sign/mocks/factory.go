package mocks

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCryptoFactory struct {
	mock.Mock
}

func (m *MockCryptoFactory) CreateMarshaller(input domain.AlgorithmType) (crypto.AlgorithmMarshaller, error) {
	args := m.Called(input)
	return args.Get(0).(crypto.AlgorithmMarshaller), args.Error(1)
}
