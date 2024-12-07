package device

import (
	"fmt"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
)

type DeviceService interface {
	GetById(id string) (*domain.Device, error)
	Save(input *domain.Device) error
	GetAll(pageNr int, pageSize int) ([]*domain.Device, int, error)
}

type DeviceRepository interface {
	Save(device domain.Device) error
	FindByID(id string) (*domain.Device, error)
	GetAll(pageNr int, pageSize int) ([]*domain.Device, int, error)
}

type CryptoFactory interface {
	CreateMarshaller(input domain.AlgorithmType) (crypto.AlgorithmMarshaller, error)
	GenerateAlgorithm(input domain.AlgorithmType) (crypto.Signer, error)
}

type SignatureDeviceServiceImpl struct {
	repository DeviceRepository
	factory    CryptoFactory
}

func NewDeviceService(repository DeviceRepository, factory CryptoFactory) *SignatureDeviceServiceImpl {
	return &SignatureDeviceServiceImpl{
		repository: repository,
		factory:    factory,
	}
}

func (s *SignatureDeviceServiceImpl) GetAll(pageNr int, pageSize int) ([]*domain.Device, int, error) {
	if pageNr < 1 || pageSize < 1 {
		return nil, 0, services.NewServiceError("invalid page number or page size", http.StatusBadRequest)
	}
	return s.repository.GetAll(pageNr, pageSize)
}

func (s *SignatureDeviceServiceImpl) GetById(id string) (*domain.Device, error) {
	device, err := s.repository.FindByID(id)
	if err != nil || device == nil {
		return nil, err
	}
	return device, nil
}

func (s *SignatureDeviceServiceImpl) Save(input *domain.Device) error {
	if input == nil {
		return services.NewServiceError(fmt.Sprintf("invalid request"), http.StatusBadRequest)

	}
	if input.ID == "" {
		return services.NewServiceError(fmt.Sprintf("id is a required field"), http.StatusBadRequest)
	}

	publicKey, privateKey, err := s.createAlgorithmForType(input.AlgorithmType)
	if err != nil {
		return err
	}
	input.PublicKey = publicKey
	input.PrivateKey = privateKey

	err = s.repository.Save(*input)
	if err != nil {
		return err
	}
	return nil
}

func (s *SignatureDeviceServiceImpl) createAlgorithmForType(algorithmType domain.AlgorithmType) ([]byte, []byte, error) {
	generatedAlgorithm, err := s.factory.GenerateAlgorithm(algorithmType)
	if err != nil {
		return nil, nil, err
	}

	marshaller, err := s.factory.CreateMarshaller(algorithmType)
	if err != nil {
		return nil, nil, err
	}
	return marshaller.Encode(generatedAlgorithm)
}
