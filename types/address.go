package types

import (
	"encoding/hex"
	"fmt"
)

// Address represents a 20-byte address.
type Address [20]byte

// ToSlice converts the Address to a byte slice.
func (a Address) ToSlice() []byte {
	return a[:]
}

// String returns the hexadecimal representation of the Address.
func (a Address) String() string {
	return hex.EncodeToString(a[:])
}

// AddressFromBytes creates an Address from a byte slice.
// It panics if the byte slice length is not 20.
func AddressFromBytes(b []byte) (Address, error) {
	if len(b) != 20 {
		return Address{}, fmt.Errorf("invalid byte length: expected 20, got %d", len(b))
	}

	var addr Address
	copy(addr[:], b)
	return addr, nil
}
