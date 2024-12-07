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
	GetDeviceCounterAndLastEncoded(id string) (int64, string, error)
	SaveDeviceCounterAndLastEncoded(id string, counter int64, currentSignature string) error
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

func (sc *SignatureCreation) SignTransaction(data []byte, deviceID string) (string, string, error) {
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

	return sc.signTransaction(device, data)
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
	err = sc.repository.SaveDeviceCounterAndLastEncoded(device.ID, counter, currentSignatureEncoded)
	if err != nil {
		sc.counterMu.Unlock()
		return "", "", err
	}
	sc.counterMu.Unlock()

	if counter == 0 {
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