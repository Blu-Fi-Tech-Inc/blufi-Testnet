package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"errors"

	"github.com/blu-fi-tech-inc/boriqua_project/types"
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

type PublicKey struct {
	*ecdsa.PublicKey
}

func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return &PrivateKey{privKey}, &PublicKey{&privKey.PublicKey}, nil
}

func (pub *PublicKey) Address() (types.Address, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub.PublicKey)
	if err != nil {
		return types.Address{}, err
	}

	hash := sha256.Sum256(pubBytes)
	addr := types.Address(hex.EncodeToString(hash[:]))
	return addr, nil
}

func (priv *PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv.PrivateKey, data)
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
}

func (pub *PublicKey) Verify(data []byte, sig *Signature) bool {
	return ecdsa.Verify(pub.PublicKey, data, sig.R, sig.S)
}

type Signature struct {
	R, S *big.Int
}
