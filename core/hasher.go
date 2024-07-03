package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"

	"github.com/blu-fi-tech-inc/boriqua_project/types"
)

// Hasher is an interface that defines a method to hash various types of data.
type Hasher interface {
	Hash(interface{}) types.Hash
}

// BlockHasher implements the Hasher interface for blocks.
type BlockHasher struct{}

// Hash computes the hash of a block header.
func (BlockHasher) Hash(b interface{}) types.Hash {
	header, ok := b.(*Header)
	if !ok {
		log.Fatalf("BlockHasher: expected *Header, got %T", b)
	}

	h := sha256.Sum256(header.Bytes())
	return types.Hash(h)
}

// TxHasher implements the Hasher interface for transactions.
type TxHasher struct{}

// Hash computes the hash of a transaction.
func (TxHasher) Hash(tx interface{}) types.Hash {
	t, ok := tx.(*Transaction)
	if !ok {
		log.Fatalf("TxHasher: expected *Transaction, got %T", tx)
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, t.Data); err != nil {
		log.Fatalf("failed to write tx.Data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, t.To); err != nil {
		log.Fatalf("failed to write tx.To: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, t.Value); err != nil {
		log.Fatalf("failed to write tx.Value: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, t.From); err != nil {
		log.Fatalf("failed to write tx.From: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, t.Nonce); err != nil {
		log.Fatalf("failed to write tx.Nonce: %v", err)
	}

	return types.Hash(sha256.Sum256(buf.Bytes()))
}
