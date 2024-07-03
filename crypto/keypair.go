package crypto

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
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

func (priv *PrivateKey) Sign(data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    r, s, err := ecdsa.Sign(rand.Reader, priv.PrivateKey, hash[:])
    if err != nil {
        return nil, err
    }
    return append(r.Bytes(), s.Bytes()...), nil
}

func (pub *PublicKey) Verify(data, signature []byte) bool {
    r := new(big.Int).SetBytes(signature[:len(signature)/2])
    s := new(big.Int).SetBytes(signature[len(signature)/2:])
    return ecdsa.Verify(pub.PublicKey, data, r, s)
}
