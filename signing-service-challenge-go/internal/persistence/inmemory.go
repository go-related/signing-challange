package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
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

func (in *InMemoryStorage) GetAll(pageNr int, pageSize int) ([]*domain.SignatureDevice, int, error) {
	startIndex := (pageNr - 1) * pageSize
	if startIndex >= len(in.signatureDevices) {
		return []*domain.SignatureDevice{}, len(in.signatureDevices), nil // No devices to return
	}
	endIndex := startIndex + pageSize
	if endIndex > len(in.signatureDevices) {
		endIndex = len(in.signatureDevices)
	}

	counter := 0
	result := make([]*domain.SignatureDevice, pageSize)
	i := 0
	for _, device := range in.signatureDevices {
		if counter > endIndex {
			break
		}

		if counter >= startIndex && i < pageSize {
			result[i] = &device
			i++
		}
		counter++
	}
	return result, len(in.signatureDevices), nil
}

func (in *InMemoryStorage) Save(device domain.SignatureDevice) error {
	if _, exists := in.signatureDevices[device.ID]; exists {
		return services.NewDBError("invalid id for the device")
	}
	in.signatureDevices[device.ID] = device
	return nil
}

func (in *InMemoryStorage) FindByID(id string) (*domain.SignatureDevice, error) {
	current, exists := in.signatureDevices[id]
	if !exists {
		return nil, services.NewDBError("invalid id for the device")
	}
	return &current, nil
}

func (in *InMemoryStorage) GetDeviceCounterAndLastEncoded(id string) (int64, string, error) {
	current, exists := in.signedTransactions[id]
	if !exists {
		return 0, "", services.NewDBError("invalid id for the device")
	}
	lastData := current[len(current)-1]
	return lastData.Counter, lastData.Signature, nil
}

func (in *InMemoryStorage) SaveDeviceCounterAndLastEncoded(id string, counter int64, currentSignature string) error {
	currentDevice, exists := in.signatureDevices[id]
	if !exists {
		return services.NewDBError("invalid id for the device")
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
