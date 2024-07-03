package core

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) IsZero() bool {
	emptyHash := Hash{}
	return h == emptyHash
}

type Address string

type Signature struct {
	R, S *big.Int
}

type PublicKey struct {
	X, Y *big.Int
}

func (pk PublicKey) Verify(hash Hash, signature Signature) bool {
	curve := elliptic.P256()
	r, s := signature.R, signature.S
	x, y := curve.ScalarBaseMult(hash[:])
	return curve.Verify(pk.X, pk.Y, x, y, r, s)
}
