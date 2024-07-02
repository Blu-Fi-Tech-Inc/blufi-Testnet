package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"

	"github.com/blu-fi-tech-inc/boriqua_project/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type TxHasher struct{}

// Hash will hash the whole bytes of the TX without exception.
func (TxHasher) Hash(tx *Transaction) types.Hash {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, tx.Data); err != nil {
		log.Fatalf("failed to write tx.Data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, tx.To); err != nil {
		log.Fatalf("failed to write tx.To: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, tx.Value); err != nil {
		log.Fatalf("failed to write tx.Value: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, tx.From); err != nil {
		log.Fatalf("failed to write tx.From: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, tx.Nonce); err != nil {
		log.Fatalf("failed to write tx.Nonce: %v", err)
	}

	return types.Hash(sha256.Sum256(buf.Bytes()))
}
