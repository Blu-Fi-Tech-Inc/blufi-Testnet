package crypto

import (
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/x509"
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
    pubBytes, err := x509.MarshalPKIXPublicKey(pub.PublicKey)
    if err != nil {
        return types.Address{}, err
    }

    hash := sha256.Sum256(pubBytes)
    addr := types.Address(hex.EncodeToString(hash[:]))
    return addr, nil
}
