package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	"github.com/google/uuid"
)

type InMemoryStorage struct {
	signatureDevices map[string]*domain.SignatureDevice
	signedCreations  map[string]*[]*domain.SignedCreations
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		signatureDevices: map[string]*domain.SignatureDevice{},
		signedCreations:  map[string]*[]*domain.SignedCreations{},
	}
}

func (in *InMemoryStorage) GetAll(pageNr int, pageSize int) ([]*domain.SignatureDevice, int, error) {
	startIndex := (pageNr - 1) * pageSize
	if startIndex >= len(in.signatureDevices) {
		return []*domain.SignatureDevice{}, len(in.signatureDevices), nil
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
			result[i] = device
			i++
		}
		counter++
	}
	return result, len(in.signatureDevices), nil
}

func (in *InMemoryStorage) GetAllSigningCreations(deviceId string, pageNr int, pageSize int) ([]*domain.SignedCreations, int, error) {
	creationsList, exist := in.signedCreations[deviceId]
	if !exist {
		return nil, 0, nil
	}
	creations := *creationsList

	startIndex := (pageNr - 1) * pageSize
	if startIndex >= len(creations) {
		return []*domain.SignedCreations{}, len(creations), services.NewDBError("invalid page number")
	}
	endIndex := startIndex + pageSize
	if endIndex > len(creations) {
		endIndex = len(creations)
	}

	counter := 0
	result := make([]*domain.SignedCreations, pageSize)
	for _, data := range creations[startIndex:endIndex] {
		result[counter] = data
		counter++
	}

	return result, len(creations), nil
}

func (in *InMemoryStorage) Save(device domain.SignatureDevice) error {
	if _, exists := in.signatureDevices[device.ID]; exists {
		return services.NewDBError("invalid id for the device")
	}
	var signingCreations []*domain.SignedCreations
	in.signatureDevices[device.ID] = &device
	in.signedCreations[device.ID] = &signingCreations
	return nil
}

func (in *InMemoryStorage) FindByID(id string) (*domain.SignatureDevice, error) {
	current, exists := in.signatureDevices[id]
	if !exists {
		return nil, services.NewDBError("invalid id for the device")
	}
	return current, nil
}

func (in *InMemoryStorage) GetDeviceCounterAndLastEncoded(id string) (int64, string, error) {
	current, exists := in.signedCreations[id]
	if !exists {
		return 0, "", services.NewDBError("invalid id for the device")
	}
	if current == nil || len(*current) == 0 {
		return 0, "", nil
	}
	lastData := (*current)[len(*current)-1]
	return lastData.Counter, lastData.Signature, nil
}

func (in *InMemoryStorage) SaveDeviceCounterAndLastEncoded(id string, counter int64, currentSignature, signedData string) error {
	currentDevice, exists := in.signatureDevices[id]
	if !exists {
		return services.NewDBError("invalid id for the device")
	}
	currentDevice.Counter = counter

	list := in.signedCreations[id]
	currentData := *list

	currentData = append(currentData, &domain.SignedCreations{
		ID:         uuid.New().String(),
		Counter:    counter,
		Signature:  currentSignature,
		SignedData: signedData,
	})
	in.signedCreations[id] = &currentData
	return nil
}
