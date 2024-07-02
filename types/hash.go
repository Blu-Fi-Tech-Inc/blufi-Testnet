package types

import (
	"encoding/hex"
	"fmt"
)

// Hash represents a 32-byte hash.
type Hash [32]byte

// IsZero checks if the Hash is all zeros.
func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// ToSlice converts the Hash to a byte slice.
func (h Hash) ToSlice() []byte {
	return h[:]
}

// String returns the hexadecimal representation of the Hash.
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// HashFromBytes creates a Hash from a byte slice.
// It returns an error if the byte slice length is not 32.
func HashFromBytes(b []byte) (Hash, error) {
	if len(b) != 32 {
		return Hash{}, fmt.Errorf("invalid byte length: expected 32, got %d", len(b))
	}

	var hash Hash
	copy(hash[:], b)
	return hash, nil
}
