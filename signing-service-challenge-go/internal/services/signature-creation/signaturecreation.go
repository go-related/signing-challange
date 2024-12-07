package signature_creation

import (
	"encoding/base64"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	"net/http"
	"sync"
)

type SignatureDeviceRepository interface {
	FindByID(id string) (*domain.SignatureDevice, error)
	GetAllSigningCreations(deviceId string, pageNr int, pageSize int) ([]*domain.SignedCreations, int, error)
	GetDeviceCounterAndLastEncoded(id string) (int64, string, error)
	SaveDeviceCounterAndLastEncoded(id string, counter int64, currentSignature, data string) error
}

type SignatureCreation struct {
	repository SignatureDeviceRepository
	counterMu  sync.Mutex
}

func NewSignatureCreation(repository SignatureDeviceRepository) *SignatureCreation {
	return &SignatureCreation{
		repository: repository,
		counterMu:  sync.Mutex{},
	}
}
func (sc *SignatureCreation) GetAllSigningCreations(deviceId string, pageNr int, pageSize int) ([]*domain.SignedCreations, int, error) {
	if deviceId == "" {
		return nil, 0, services.NewServiceError("deviceId is required", http.StatusBadRequest)
	}
	if pageNr <= 0 {
		return nil, 0, services.NewServiceError("pageNr is required", http.StatusBadRequest)
	}
	if pageSize <= 0 {
		return nil, 0, services.NewServiceError("pageSize is required", http.StatusBadRequest)
	}

	return sc.repository.GetAllSigningCreations(deviceId, pageNr, pageSize)
}

func (sc *SignatureCreation) SignTransaction(deviceID string, data []byte) (string, string, error) {
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

func (sc *SignatureCreation) signTransaction(device *domain.SignatureDevice, data []byte) (string, string, error) {
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

func (sc *SignatureCreation) loadKeyFromDevice(device *domain.SignatureDevice) (crypto.Signer, error) {
	switch device.AlgorithmType {
	case domain.AlgorithmTypeRSA:
		marshaller := crypto.NewRSAMarshaler()
		key, err := marshaller.Unmarshal(device.PrivateKey)
		if err != nil {
			return nil, err
		}
		return key, nil
	case domain.AlgorithmTypeECC:
		marshaller := crypto.NewECCMarshaler()
		key, err := marshaller.Decode(device.PrivateKey)
		if err != nil {
			return nil, err
		}
		return key, nil
	default:
		return nil, services.NewServiceError("invalid algorithm type registered for this device", http.StatusBadRequest)
	}
}
