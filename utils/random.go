package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
	"github.com/blu-fi-tech-inc/boriqua_project/crypto"
	"github.com/blu-fi-tech-inc/boriqua_project/types"
	"github.com/stretchr/testify/assert"
)

// RandomBytes generates a slice of random bytes of given size.
func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

// RandomHash generates a random Hash.
func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

// NewRandomTransaction creates a new random transaction without signature.
func NewRandomTransaction(size int) *core.Transaction {
	return core.NewTransaction(types.Address(RandomHash()), types.Address(RandomHash()), uint64(rand.Intn(1000)), RandomBytes(size))
}

// NewRandomTransactionWithSignature creates a new random transaction and signs it with the provided private key.
func NewRandomTransactionWithSignature(t *testing.T, privKey crypto.PrivateKey, size int) *core.Transaction {
	tx := NewRandomTransaction(size)
	assert.Nil(t, tx.Sign(privKey))
	return tx
}

// NewRandomBlock creates a new random block with a single random signed transaction.
func NewRandomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *core.Block {
	txSigner := crypto.GeneratePrivateKey()
	tx := NewRandomTransactionWithSignature(t, txSigner, 100)
	header := &core.Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	b, err := core.NewBlock(header, []*core.Transaction{tx})
	assert.Nil(t, err)

	dataHash, err := core.CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash

	return b
}

// NewRandomBlockWithSignature creates a new random block and signs it with the provided private key.
func NewRandomBlockWithSignature(t *testing.T, pk crypto.PrivateKey, height uint32, prevHash types.Hash) *core.Block {
	b := NewRandomBlock(t, height, prevHash)
	assert.Nil(t, b.Sign(pk))
	return b
}
