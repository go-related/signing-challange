package signature_device

import (
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	"net/http"
)

type SignatureDeviceRepository interface {
	Save(device domain.SignatureDevice) error
	FindByID(id string) (*domain.SignatureDevice, error)
	GetAll(pageNr int, pageSize int) ([]*domain.SignatureDevice, int, error)
}

type SignatureDeviceServiceImpl struct {
	repository SignatureDeviceRepository
}

func NewDeviceService(repository SignatureDeviceRepository) *SignatureDeviceServiceImpl {
	return &SignatureDeviceServiceImpl{
		repository: repository,
	}
}

func (s *SignatureDeviceServiceImpl) GetAll(pageNr int, pageSize int) ([]*domain.SignatureDevice, int, error) {
	if pageNr < 1 || pageSize < 1 {
		return nil, 0, services.NewServiceError("invalid page number or page size", http.StatusBadRequest)
	}
	return s.repository.GetAll(pageNr, pageSize)
}

func (s *SignatureDeviceServiceImpl) GetById(id string) (*domain.SignatureDevice, error) {
	device, err := s.repository.FindByID(id)
	if err != nil || device == nil {
		return nil, err
	}
	return device, nil
}

func (s *SignatureDeviceServiceImpl) Save(input *domain.SignatureDevice) error {
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
	switch algorithmType {
	case domain.AlgorithmTypeECC:
		generator := crypto.ECCGenerator{}
		key, err := generator.Generate()
		if err != nil {
			return nil, nil, err
		}
		marshaller := crypto.NewECCMarshaler()
		pbKey, pvKey, err := marshaller.Encode(key)
		return pbKey, pvKey, err
	case domain.AlgorithmTypeRSA:
		generator := crypto.RSAGenerator{}
		key, err := generator.Generate()
		if err != nil {
			return nil, nil, err
		}
		marshaller := crypto.NewRSAMarshaler()
		pbKey, pvKey, err := marshaller.Marshal(key)
		return pbKey, pvKey, err
	default:
		return nil, nil, services.NewServiceError(fmt.Sprintf("unknown algorithm type :%s", algorithmType), http.StatusBadRequest)
	}
}
