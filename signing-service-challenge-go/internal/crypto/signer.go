package crypto

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	// VerifySignature(data []byte, signature []byte) error // not sure if we will need it here but, it's a convenience to have it for sure
}
