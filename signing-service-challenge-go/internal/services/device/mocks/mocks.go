package mocks

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) Save(device domain.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockDeviceRepository) FindByID(id string) (*domain.Device, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *MockDeviceRepository) GetAll(pageNr int, pageSize int) ([]*domain.Device, int, error) {
	args := m.Called(pageNr, pageSize)
	return args.Get(0).([]*domain.Device), args.Int(1), args.Error(2)
}
