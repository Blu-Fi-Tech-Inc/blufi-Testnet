package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// PrivateKey represents a private key for ECDSA.
type PrivateKey struct {
	*ecdsa.PrivateKey
}

// PublicKey represents a public key for ECDSA.
type PublicKey struct {
	*ecdsa.PublicKey
}

// GenerateKeyPair generates a new ECDSA key pair.
func GenerateKeyPair() (PrivateKey, PublicKey, error) {
	curve := elliptic.P256()
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return PrivateKey{}, PublicKey{}, fmt.Errorf("error generating ECDSA key pair: %v", err)
	}

	return PrivateKey{priv}, PublicKey{&priv.PublicKey}, nil
}

// Sign signs a message using ECDSA.
func (privKey PrivateKey) Sign(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, privKey.PrivateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("error signing message: %v", err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// VerifySignature verifies a signature against a message using ECDSA.
func (pubKey PublicKey) VerifySignature(message []byte, signature []byte) (bool, error) {
	hash := sha256.Sum256(message)
	rBytes := signature[:len(signature)/2]
	sBytes := signature[len(signature)/2:]

	r := big.NewInt(0).SetBytes(rBytes)
	s := big.NewInt(0).SetBytes(sBytes)

	valid := ecdsa.Verify(pubKey.PublicKey, hash[:], r, s)
	return valid, nil
}
