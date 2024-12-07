package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
)

// ECCKeyPair is a DTO that holds ECC private and public keys.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

func (kp *ECCKeyPair) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, kp.Private, hash[:])
	if err != nil {
		return nil, fmt.Errorf("signing failed: %v", err)
	}

	// Marshal signature into ASN.1 format
	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{
		R: r,
		S: s,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature: %v", err)
	}

	return signature, nil
}

func (kp *ECCKeyPair) VerifySignature(data []byte, signature []byte) error {
	// Hash the original data
	hash := sha256.Sum256(data)

	// Unmarshal the signature
	var sigStruct struct {
		R, S *big.Int
	}
	_, err := asn1.Unmarshal(signature, &sigStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal signature: %v", err)
	}

	// Verify the signature
	if !ecdsa.Verify(kp.Public, hash[:], sigStruct.R, sigStruct.S) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// ECCMarshaler can encode and decode an ECC key pair.
type ECCMarshaler struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCMarshaler() *ECCMarshaler {
	return &ECCMarshaler{}
}

// Encode takes an ECCKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m *ECCMarshaler) Encode(keyPair Signer) ([]byte, []byte, error) {
	input := keyPair.(*ECCKeyPair)
	privateKeyBytes, err := x509.MarshalECPrivateKey(input.Private)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(input.Public)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Decode assembles an ECCKeyPair from an encoded private key.
func (m *ECCMarshaler) Decode(privateKeyBytes []byte) (Signer, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
