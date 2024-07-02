package utils

import (
	"crypto/sha256"
	"errors"
)

type Hash [32]byte

func HashFromBytes(data []byte) (Hash, error) {
	if len(data) == 0 {
		return Hash{}, errors.New("data is empty")
	}
	hash := sha256.Sum256(data)
	return hash, nil
}
