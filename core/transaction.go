package core

import (
	"crypto/ecdsa"
	"encoding/gob"
	"fmt"
	"math/rand"

	"github.com/blu-fi-tech-inc/blufi-network/crypto"
	"github.com/blu-fi-tech-inc/blufi-network/types"
)

type TxType byte

const (
	TxTypeCollection TxType = iota // 0x0
	TxTypeMint                     // 0x01
)

type CollectionTx struct {
	Fee      int64
	MetaData []byte
}

type MintTx struct {
	Fee             int64
	NFT             types.Hash
	Collection      types.Hash
	MetaData        []byte
	CollectionOwner crypto.PublicKey
	Signature       []byte
}

type Transaction struct {
	TxInner   interface{}       // Generic type for handling inner transactions
	Data      []byte
	To        crypto.PublicKey
	Value     uint64
	From      crypto.PublicKey
	Signature []byte
	Nonce     int64

	// Cached version of the tx data hash
	hash types.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:  data,
		Nonce: rand.Int63n(1000000000000000),
	}
}

func (tx *Transaction) Hash(hasher Hasher) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

func (tx *Transaction) Sign(privKey *crypto.PrivateKey) error {
	hash := tx.Hash(TxHasher{})
	sig, err := privKey.Sign(hash[:]) // Corrected to use hash as byte slice
	if err != nil {
		return err
	}

	tx.From = crypto.PublicKey{PublicKey: &privKey.PublicKey}
    tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	hash := tx.Hash(TxHasher{})
    pubKey := &tx.From
    if !crypto.VerifySignature(pubKey, hash[:], tx.Signature) {
        return fmt.Errorf("invalid transaction signature")
    }

	// Verify the inner transaction if exists
	switch innerTx := tx.TxInner.(type) {
	case CollectionTx:
		// Add specific verification for CollectionTx if needed
	case MintTx:
        if !crypto.VerifySignature(&innerTx.CollectionOwner, innerTx.Collection[:], innerTx.Signature) {
            return fmt.Errorf("invalid mint transaction signature")
        }
    }

	return nil
}

func (tx *Transaction) Decode(dec *gob.Decoder) error {
    return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc *gob.Encoder) error {
    return enc.Encode(tx)
}

func init() {
    gob.Register(CollectionTx{})
    gob.Register(MintTx{})
}
