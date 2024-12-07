package domain

import (
	"time"
)

// AlgorithmType will contain all algorithm types that this system wil support
type AlgorithmType string

var (
	AlgorithmTypeUnknown AlgorithmType = ""
	AlgorithmTypeECC     AlgorithmType = "ECC"
	AlgorithmTypeRSA     AlgorithmType = "RSA"
)

type SignatureDevice struct {
	ID            string
	AlgorithmType AlgorithmType
	Label         *string
	Counter       int64

	PublicKey  []byte //storing public key is not needed actually
	PrivateKey []byte
}

// DeviceSigner would be in case we want for one device to be able to work with multiple signers and select one of them to work each time you want to sign something.
// since there is no mentioning of this possibility on the topics I am gonna keep thing simple and store the keys in the device itself
type DeviceSigner struct {
	ID        string
	DeviceId  string
	isValid   bool
	CreatedAt time.Time

	AlgorithmType AlgorithmType
	PublicKey     []byte
	PrivateKey    []byte
}

// SignedCreations I am not sure if we need this at this time, but for historic reasons i am leaving it here
type SignedCreations struct {
	ID         string
	DeviceId   string
	Counter    int64
	Signature  string
	SignedData string
}
