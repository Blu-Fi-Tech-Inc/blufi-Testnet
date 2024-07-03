package crypto

import (
    "crypto/ecdsa"
    "crypto/rand"
    "math/big"
)

type PrivateKey struct {
    *ecdsa.PrivateKey
}

type PublicKey struct {
    *ecdsa.PublicKey
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
