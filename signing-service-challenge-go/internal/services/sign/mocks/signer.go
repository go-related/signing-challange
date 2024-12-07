package mocks

import "github.com/stretchr/testify/mock"

type MockSigner struct {
	mock.Mock
}

func (m *MockSigner) Sign(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}
