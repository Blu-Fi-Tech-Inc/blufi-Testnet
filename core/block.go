package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/blu-fi-tech-inc/blufi-network/crypto"
	"github.com/blu-fi-tech-inc/blufi-network/types"
)

// Header represents the header of a block.
type Header struct {
	Version       uint32       // Version of the block
	DataHash      types.Hash   // Hash of the block's data
	PrevBlockHash types.Hash   // Hash of the previous block's header
	Height        uint32       // Height of the block in the blockchain
	Timestamp     int64        // Timestamp when the block was created
}

// Bytes serializes the header into a byte slice using gob encoding.
func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(h); err != nil {
		panic(err) // Handle the error appropriately in production code
	}
	return buf.Bytes()
}

// Block represents a block in the blockchain.
type Block struct {
	*Header             // Pointer to the block's header
	Transactions []*Transaction // List of transactions in the block
	Validator    crypto.PublicKey // Public key of the validator who created the block
	Signature    []byte
    hash         types.Hash
}

// NewBlock creates a new block with the given header and transactions.
func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txx,
	}, nil
}

// NewBlockFromPrevHeader creates a new block based on the previous block's header and given transactions.
func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, txx)
}

// AddTransaction adds a transaction to the block and recalculates the block's data hash.
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
	hash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		// Handle the error appropriately in production code
		return
	}
	b.DataHash = hash
}

// Update the Sign method in Block struct
func (b *Block) Sign(privKey *crypto.PrivateKey) error {
	hash := b.Hash(BlockHasher{})
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = crypto.PublicKey{PublicKey: &privKey.PublicKey}
	b.Signature = sig

	return nil
}

// Verify verifies the integrity and validity of the block.
func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !crypto.VerifySignature(&b.Validator, b.Header.Bytes(), b.Signature) {
        return fmt.Errorf("block has an invalid signature")
	}

	for _, tx := range b.Transactions {
        if err := tx.Verify(); err != nil {
            return err
        }
	}

	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}

	return nil
}

// Decode decodes the block using the provided decoder.
func (b *Block) Decode(dec *gob.Decoder) error {
    return dec.Decode(b)
}

// Encode encodes the block using the provided encoder.
func (b *Block) Encode(enc *gob.Encoder) error {
    return enc.Encode(b)
}

// Hash computes and returns the hash of the block's header using the provided hasher.
func (b *Block) Hash(hasher Hasher) types.Hash {
    if b.hash.IsZero() {
        b.hash = hasher.Hash(b.Header)
    }

    return b.hash
}

// CalculateDataHash computes the hash of the block's data (transactions).
func CalculateDataHash(txx []*Transaction) (types.Hash, error) {
    buf := &bytes.Buffer{}

    for _, tx := range txx {
        if err := tx.Encode(gob.NewEncoder(buf)); err != nil {
            return types.Hash{}, err
        }
    }

    return types.Hash(sha256.Sum256(buf.Bytes())), nil
}
