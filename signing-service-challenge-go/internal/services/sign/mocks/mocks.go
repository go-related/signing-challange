package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
)

type MockSignRepository struct {
	mock.Mock
}

func (m *MockSignRepository) FindByID(id string) (*domain.Device, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *MockSignRepository) GetAllSignings(deviceId string, pageNr int, pageSize int) ([]*domain.Signings, int, error) {
	args := m.Called(deviceId, pageNr, pageSize)
	return args.Get(0).([]*domain.Signings), args.Int(1), args.Error(2)
}

func (m *MockSignRepository) GetDeviceCounterAndLastEncoded(id string) (int64, string, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.String(1), args.Error(2)
}

func (m *MockSignRepository) SaveDeviceCounterAndLastEncoded(id string, counter int64, currentSignature, data string) error {
	args := m.Called(id, counter, currentSignature, data)
	return args.Error(0)
}
