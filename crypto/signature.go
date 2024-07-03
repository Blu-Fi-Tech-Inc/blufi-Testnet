package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

// PrivateKey represents an ECDSA private key.
type PrivateKey struct {
	*ecdsa.PrivateKey
}

// PublicKey represents an ECDSA public key.
type PublicKey struct {
	*ecdsa.PublicKey
}

// GenerateKeyPair generates a new ECDSA key pair.
func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(ecdsa.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return &PrivateKey{privKey}, &PublicKey{&privKey.PublicKey}, nil
}

// Sign creates a signature for the given data using the private key.
func (privKey *PrivateKey) Sign(data []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey.PrivateKey, dataHash(data))
	if err != nil {
		return nil, err
	}
	// Encode r and s into a single byte slice
	signature, err := encodeSignature(r, s)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Verify checks if the given signature is valid for the provided data and public key.
func (pubKey *PublicKey) Verify(data []byte, signature []byte) bool {
	r, s, err := decodeSignature(signature)
	if err != nil {
		return false
	}
	return ecdsa.Verify(pubKey.PublicKey, dataHash(data), r, s)
}

// Helper function to hash data using SHA-256
func dataHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// Helper function to encode ECDSA signature (r, s) into a byte slice
func encodeSignature(r, s *big.Int) ([]byte, error) {
	const sigLen = 32 // bytes per big.Int in ECDSA signature
	sig := make([]byte, 2*sigLen)
	rBytes, sBytes := r.Bytes(), s.Bytes()
	copy(sig[sigLen-len(rBytes):sigLen], rBytes)
	copy(sig[2*sigLen-len(sBytes):2*sigLen], sBytes)
	return sig, nil
}

// Helper function to decode ECDSA signature byte slice into (r, s)
func decodeSignature(sig []byte) (*big.Int, *big.Int, error) {
	const sigLen = 32 // bytes per big.Int in ECDSA signature
	if len(sig) != 2*sigLen {
		return nil, nil, errors.New("invalid signature length")
	}
	r := new(big.Int).SetBytes(sig[:sigLen])
	s := new(big.Int).SetBytes(sig[sigLen:])
	return r, s, nil
}

// AddressFromPublicKey generates an address from a public key.
func AddressFromPublicKey(pubKey *PublicKey) string {
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	hash := sha256.Sum256(pubBytes)
	return hex.EncodeToString(hash[:20]) // Take first 20 bytes for address
}
