package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// Hash represents a 32-byte hash.
type Hash [32]byte

// String returns the hexadecimal representation of the Hash.
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// IsZero checks if the Hash is all zeros.
func (h Hash) IsZero() bool {
	emptyHash := Hash{}
	return h == emptyHash
}

// HashBytes returns the SHA-256 hash of the input bytes.
func HashBytes(data []byte) Hash {
	return sha256.Sum256(data)
}

// Address represents a 20-byte address as a string.
type Address string

// Signature represents an ECDSA signature with R and S values.
type Signature struct {
	R, S *big.Int
}

// PublicKey represents an ECDSA public key with X and Y coordinates.
type PublicKey ecdsa.PublicKey

// Verify verifies an ECDSA signature given the hash and public key.
func (pk PublicKey) Verify(hash Hash, signature Signature) bool {
	curve := elliptic.P256()
	return ecdsa.Verify((*ecdsa.PublicKey)(&pk), hash[:], signature.R, signature.S)
}
