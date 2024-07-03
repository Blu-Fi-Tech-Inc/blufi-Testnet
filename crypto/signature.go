package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

// Signature represents an ECDSA signature with R and S values.
type Signature struct {
	R, S *big.Int
}

// Sign computes the ECDSA signature of the provided data using the given private key.
func Sign(privKey *PrivateKey, data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey.PrivateKey, data)
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
}

// VerifyECDSASignature verifies an ECDSA signature given the public key, data, and signature.
func VerifyECDSASignature(pubKey *PublicKey, data []byte, signature *Signature) bool {
	return ecdsa.Verify(pubKey.PublicKey, data, signature.R, signature.S)
}
