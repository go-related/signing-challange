package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
)

type MockMarshaller struct {
	mock.Mock
}

func (m *MockMarshaller) Encode(input crypto.Signer) ([]byte, []byte, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]byte), args.Get(1).([]byte), args.Error(2)
}

func (m *MockMarshaller) Decode(input []byte) (crypto.Signer, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(crypto.Signer), args.Error(1)
}
