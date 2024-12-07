package signature_device

import "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"

type SignatureDeviceService interface {
	GetById(id string) (*domain.SignatureDevice, error)
	Save(input *domain.SignatureDevice) error
	GetAll(pageNr int, pageSize int) ([]*domain.SignatureDevice, int, error)
}
