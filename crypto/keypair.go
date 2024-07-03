package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/blu-fi-tech-inc/boriqua_project/types"
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

type PublicKey struct {
	*ecdsa.PublicKey
}

func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(ecdsa.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return &PrivateKey{privKey}, &PublicKey{&privKey.PublicKey}, nil
}

func (pub *PublicKey) Address() (types.Address, error) {
	pubBytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	hash := sha256.Sum256(pubBytes)
	addr := types.Address(hex.EncodeToString(hash[:]))
	return addr, nil
}

func (priv *PrivateKey) Sign(data []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv.PrivateKey, data)
	if err != nil {
		return nil, err
	}
	// Concatenate r and s into a single byte slice
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func (pub *PublicKey) Verify(data, signature []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)

	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	return ecdsa.Verify(pub.PublicKey, data, &r, &s)
}
