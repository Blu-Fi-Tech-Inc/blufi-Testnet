package consensus

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
)

const (
	targetBits = 24
	maxNonce   = ^uint32(0) // maximum value for a 32-bit integer
)

// ProofOfWork represents a Proof of Work consensus mechanism.
type ProofOfWork struct {
	block  *core.Block
	target *big.Int
}

// NewProofOfWork initializes a new ProofOfWork for a given block.
func NewProofOfWork(block *core.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{block, target}
}

// prepareData prepares the data for mining.
func (pow *ProofOfWork) prepareData(nonce uint32) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.Header.PrevBlockHash.ToSlice(),
			pow.block.Header.DataHash.ToSlice(),
			core.UintToBytes(pow.block.Header.Timestamp),
			core.UintToBytes(uint64(targetBits)),
			core.UintToBytes(uint64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run performs the Proof of Work algorithm.
func (pow *ProofOfWork) Run() (uint32, []byte) {
	var (
		hashInt big.Int
		hash    [32]byte
		nonce   uint32
	)

	fmt.Printf("Mining a new block")

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println("\n\n")

	return nonce, hash[:]
}

// Validate validates the Proof of Work for a block.
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Header.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

// MineBlock mines a new block using Proof of Work.
func MineBlock(prevBlock *core.Block, transactions []*core.Transaction) (*core.Block, error) {
	var (
		block *core.Block
		err   error
	)

	block = &core.Block{}
	block.Header = &core.Header{
		Version:       1,
		PrevBlockHash: prevBlock.Header.BlockHash(),
		Height:        prevBlock.Header.Height + 1,
		Timestamp:     time.Now().UnixNano(),
	}

	block.Transactions = transactions
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Header.Nonce = nonce
	block.Header.Hash = hash

	if pow.Validate() {
		return block, nil
	} else {
		err = fmt.Errorf("proof of work validation failed")
	}

	return nil, err
}
