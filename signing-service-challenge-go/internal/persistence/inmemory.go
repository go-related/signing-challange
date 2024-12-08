package persistence

import (
	"github.com/google/uuid"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
)

type InMemoryStorage struct {
	devicesData  map[string]*domain.Device
	signingsData map[string]*[]*domain.Signings
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		devicesData:  map[string]*domain.Device{},
		signingsData: map[string]*[]*domain.Signings{},
	}
}

func (in *InMemoryStorage) GetAll(pageNr int, pageSize int) ([]*domain.Device, int, error) {
	startIndex := (pageNr - 1) * pageSize
	if startIndex >= len(in.devicesData) {
		return []*domain.Device{}, len(in.devicesData), nil
	}
	endIndex := startIndex + pageSize
	if endIndex > len(in.devicesData) {
		endIndex = len(in.devicesData)
	}

	counter := 0
	result := make([]*domain.Device, pageSize)
	i := 0
	for _, device := range in.devicesData {
		if counter > endIndex {
			break
		}

		if counter >= startIndex && i < pageSize {
			result[i] = device
			i++
		}
		counter++
	}
	return result, len(in.devicesData), nil
}

func (in *InMemoryStorage) GetAllSignings(deviceId string, pageNr int, pageSize int) ([]*domain.Signings, int, error) {
	creationsList, exist := in.signingsData[deviceId]
	if !exist || creationsList == nil || len(*creationsList) == 0 {
		return nil, 0, nil
	}
	creations := *creationsList

	startIndex := (pageNr - 1) * pageSize
	if startIndex >= len(creations) {
		return []*domain.Signings{}, len(creations), services.NewDBError("invalid page number")
	}
	endIndex := startIndex + pageSize
	if endIndex > len(creations) {
		endIndex = len(creations)
	}

	counter := 0
	result := make([]*domain.Signings, pageSize)
	for _, data := range creations[startIndex:endIndex] {
		result[counter] = data
		counter++
	}

	return result, len(creations), nil
}

func (in *InMemoryStorage) Save(device domain.Device) error {
	if _, exists := in.devicesData[device.ID]; exists {
		return services.NewDBError("invalid id for the device")
	}
	var signingCreations []*domain.Signings
	in.devicesData[device.ID] = &device
	in.signingsData[device.ID] = &signingCreations
	return nil
}

func (in *InMemoryStorage) FindByID(id string) (*domain.Device, error) {
	current, exists := in.devicesData[id]
	if !exists {
		return nil, services.NewDBError("invalid id for the device")
	}
	return current, nil
}

func (in *InMemoryStorage) GetDeviceCounterAndLastEncoded(id string) (int64, string, error) {
	current, exists := in.signingsData[id]
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
	currentDevice, exists := in.devicesData[id]
	if !exists {
		return services.NewDBError("invalid id for the device")
	}
	currentDevice.Counter = counter

	list := in.signingsData[id]
	currentData := *list

	currentData = append(currentData, &domain.Signings{
		ID:         uuid.New().String(),
		Counter:    counter,
		Signature:  currentSignature,
		SignedData: signedData,
	})
	in.signingsData[id] = &currentData
	return nil
}
