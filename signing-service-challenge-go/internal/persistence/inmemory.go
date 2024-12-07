package persistence

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/google/uuid"
)

type InMemoryStorage struct {
	signatureDevices   map[string]domain.SignatureDevice
	signedTransactions map[string][]domain.SignedTransactions
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		signatureDevices:   map[string]domain.SignatureDevice{},
		signedTransactions: map[string][]domain.SignedTransactions{},
	}
}

func (in *InMemoryStorage) Save(device domain.SignatureDevice) error {
	if _, exists := in.signatureDevices[device.ID]; exists {
		return errors.New("invalid id for the device")
	}
	in.signatureDevices[device.ID] = device
	return nil
}

func (in *InMemoryStorage) FindByID(id string) (*domain.SignatureDevice, error) {
	current, exists := in.signatureDevices[id]
	if !exists {
		return nil, errors.New("invalid id for the device")
	}
	return &current, nil
}

func (in *InMemoryStorage) GetDeviceCounterAndLastEncoded(id string) (int64, string, error) {
	current, exists := in.signedTransactions[id]
	if !exists {
		return 0, "", errors.New("invalid id for the device")
	}
	lastData := current[len(current)-1]
	return lastData.Counter, lastData.Signature, nil
}

func (in *InMemoryStorage) UpdateDeviceCounterAndLastEncoded(id string, counter int64, currentSignature string) error {
	currentDevice, exists := in.signatureDevices[id]
	if !exists {
		return errors.New("invalid id for the device")
	}
	currentDevice.Counter = counter

	list := in.signedTransactions[id]

	list = append(list, domain.SignedTransactions{
		ID:        uuid.New().String(),
		Counter:   counter,
		Signature: currentSignature,
	})
	in.signedTransactions[id] = list
	return nil
}
