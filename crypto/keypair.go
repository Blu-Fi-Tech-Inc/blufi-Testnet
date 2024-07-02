package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509" // Importing x509 package
	"encoding/pem"
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
	addr, err := types.AddressFromBytes(pubBytes)
	if err != nil {
		return types.Address{}, err
	}
	return addr, nil
}
