package crypto

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
)

type AlgorithmMarshaller interface {
	Encode(input Signer) ([]byte, []byte, error)
	Decode(input []byte) (Signer, error)
}

func CreateMarshaller(input domain.AlgorithmType) (AlgorithmMarshaller, error) {
	switch input {
	case domain.AlgorithmTypeECC:
		return NewECCMarshaler(), nil
	case domain.AlgorithmTypeRSA:
		return NewRSAMarshaler(), nil
	default:
		return nil, errors.New("unknown algorithm type")
	}
}

func GenerateAlgorithm(input domain.AlgorithmType) (Signer, error) {
	switch input {
	case domain.AlgorithmTypeECC:
		var generator ECCGenerator
		return generator.Generate()
	case domain.AlgorithmTypeRSA:
		var generator RSAGenerator
		return generator.Generate()
	default:
		return nil, errors.New("unknown algorithm type")
	}
}
