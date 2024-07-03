package crypto

import (
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/x509"
    "crypto/elliptic"
    "crypto/sha256"
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
    return types.AddressFromBytes(hash[:20])
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

func VerifySignature(pub *PublicKey, data, signature []byte) bool {
    if len(signature) != 64 {
        return false
    }
    r := new(big.Int).SetBytes(signature[:32])
    s := new(big.Int).SetBytes(signature[32:])
    return ecdsa.Verify(pub.PublicKey, data, r, s)
}
