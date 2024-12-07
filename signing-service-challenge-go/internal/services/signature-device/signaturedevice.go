package signature_device

import (
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	"github.com/google/uuid"
	"net/http"
)

type SignatureDeviceRepository interface {
	Save(device domain.SignatureDevice) error
	FindByID(id string) (*domain.SignatureDevice, error)
}

type SignatureDeviceService struct {
	repository SignatureDeviceRepository
}

func NewDeviceService(repository SignatureDeviceRepository) *SignatureDeviceService {
	return &SignatureDeviceService{
		repository: repository,
	}
}

func (s *SignatureDeviceService) GetById(id string) (*domain.SignatureDevice, error) {
	device, err := s.repository.FindByID(id)
	if err != nil || device == nil {
		// TODO wrap this
		return nil, err
	}
	return device, nil
}

func (s *SignatureDeviceService) Save(algorithmType string, label *string) error {
	data := domain.SignatureDevice{
		ID:    uuid.New().String(),
		Label: label,
	}

	algType, publicKey, privateKey, err := s.createAlgorithmForType(algorithmType)
	if err != nil {
		return err
	}
	if algType == domain.AlgorithmTypeUnknown {
		return services.NewServiceError(fmt.Sprintf("unknown algorithm type: %s", algorithmType), http.StatusBadRequest)
	}
	data.AlgorithmType = algType
	data.PublicKey = publicKey
	data.PrivateKey = privateKey

	err = s.repository.Save(data)
	if err != nil {
		// TODO wrap this error
		return err
	}
	return nil
}

func (s *SignatureDeviceService) createAlgorithmForType(algorithmType string) (domain.AlgorithmType, []byte, []byte, error) {
	switch algorithmType {
	case "ECC":
		generator := crypto.ECCGenerator{}
		key, err := generator.Generate()
		if err != nil {
			return domain.AlgorithmTypeUnknown, nil, nil, err
		}
		marshaller := crypto.NewECCMarshaler()
		pbKey, pvKey, err := marshaller.Encode(key)
		return domain.AlgorithmTypeECC, pbKey, pvKey, err
	case "RSA":
		generator := crypto.RSAGenerator{}
		key, err := generator.Generate()
		if err != nil {
			return domain.AlgorithmTypeUnknown, nil, nil, err
		}
		marshaller := crypto.NewRSAMarshaler()
		pbKey, pvKey, err := marshaller.Marshal(key)
		return domain.AlgorithmTypeRSA, pbKey, pvKey, err
	default:
		return domain.AlgorithmTypeUnknown, nil, nil, services.NewServiceError(fmt.Sprintf("unknown algorithm type :%s", algorithmType), http.StatusBadRequest)
	}
}
