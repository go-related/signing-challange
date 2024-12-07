package signature_creation

type SignatureService interface {
	SignTransaction(data []byte, deviceID string) (string, string, error)
}
