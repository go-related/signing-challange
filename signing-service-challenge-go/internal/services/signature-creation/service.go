package signature_creation

import "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"

type SignatureService interface {
	SignTransaction(deviceID string, data []byte) (string, string, error)
	GetAllSigningCreations(deviceId string, pageNr int, pageSize int) ([]*domain.SignedCreations, int, error)
}
