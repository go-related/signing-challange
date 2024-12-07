package sign

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
)

type SignService interface {
	Sign(deviceID string, data []byte) (string, string, error)
	GetAllSignings(deviceId string, pageNr int, pageSize int) ([]*domain.Signings, int, error)
}

type SignRepository interface {
	FindByID(id string) (*domain.Device, error)
	GetAllSignings(deviceId string, pageNr int, pageSize int) ([]*domain.Signings, int, error)
	GetDeviceCounterAndLastEncoded(id string) (int64, string, error)
	SaveDeviceCounterAndLastEncoded(id string, counter int64, currentSignature, data string) error
}

type CryptoFactory interface {
	CreateMarshaller(input domain.AlgorithmType) (crypto.AlgorithmMarshaller, error)
	GenerateAlgorithm(input domain.AlgorithmType) (crypto.Signer, error)
}

type SignServiceImpl struct {
	repository    SignRepository
	cryptoFactory CryptoFactory
	counterMu     sync.Mutex
}

func NewSignService(repository SignRepository, factory CryptoFactory) *SignServiceImpl {
	return &SignServiceImpl{
		repository:    repository,
		cryptoFactory: factory,
		counterMu:     sync.Mutex{},
	}
}
func (sc *SignServiceImpl) GetAllSignings(deviceId string, pageNr int, pageSize int) ([]*domain.Signings, int, error) {
	if deviceId == "" {
		return nil, 0, services.NewServiceError("deviceId is required", http.StatusBadRequest)
	}
	if pageNr <= 0 {
		return nil, 0, services.NewServiceError("pageNr is required", http.StatusBadRequest)
	}
	if pageSize <= 0 {
		return nil, 0, services.NewServiceError("pageSize is required", http.StatusBadRequest)
	}

	return sc.repository.GetAllSignings(deviceId, pageNr, pageSize)
}

func (sc *SignServiceImpl) Sign(deviceID string, data []byte) (string, string, error) {
	if deviceID == "" {
		return "", "", services.NewServiceError("device_id is a required field", http.StatusBadRequest)
	}

	if len(data) == 0 {
		return "", "", services.NewServiceError("data is a required field", http.StatusBadRequest)
	}

	device, err := sc.repository.FindByID(deviceID)
	if err != nil {
		return "", "", err
	}
	if device == nil {
		return "", "", services.NewServiceError("invalid device_id value", http.StatusBadRequest)
	}

	signature, signedData, err := sc.signTransaction(device, data)
	if err != nil {
		return "", "", err
	}
	return signature, signedData, nil
}

func (sc *SignServiceImpl) signTransaction(device *domain.Device, data []byte) (string, string, error) {
	signer, err := sc.loadKeyFromDevice(device)
	if err != nil {
		return "", "", err
	}

	signature, err := signer.Sign(data)
	if err != nil {
		return "", "", err
	}

	currentSignatureEncoded := base64.StdEncoding.EncodeToString(signature)

	sc.counterMu.Lock()
	counter, lastEncoded, err := sc.repository.GetDeviceCounterAndLastEncoded(device.ID)
	if err != nil {
		sc.counterMu.Unlock()
		return "", "", err
	}
	counter += 1
	err = sc.repository.SaveDeviceCounterAndLastEncoded(device.ID, counter, currentSignatureEncoded, string(data))
	if err != nil {
		sc.counterMu.Unlock()
		return "", "", err
	}
	sc.counterMu.Unlock()

	if counter == 1 {
		lastEncoded = base64.StdEncoding.EncodeToString([]byte(device.ID))
	}
	signedData := fmt.Sprintf("%d_%s_%s", counter, string(data), lastEncoded)

	return currentSignatureEncoded, signedData, nil
}

func (sc *SignServiceImpl) loadKeyFromDevice(device *domain.Device) (crypto.Signer, error) {
	marshaller, err := sc.cryptoFactory.CreateMarshaller(device.AlgorithmType)
	if err != nil {
		return nil, err
	}
	key, err := marshaller.Decode(device.PrivateKey)
	if err != nil {
		return nil, err
	}
	return key, nil
}
