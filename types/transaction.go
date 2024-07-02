package types

import (
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	From   Address
	To     Address
	Amount uint64
	Data   []byte
}

func NewTransaction(from, to Address, amount uint64, data []byte) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
		Data:   data,
	}
}

func (tx *Transaction) Hash() Hash {
	// Example hashing logic (you may use a more complex hashing scheme)
	hash, err := HashFromBytes([]byte(fmt.Sprintf("%v%v%v%v", tx.From, tx.To, tx.Amount, tx.Data)))
	if err != nil {
		log.Fatalf("Failed to get hash from bytes: %v", err)
	}
	return hash
}

func (tx *Transaction) String() string {
	return fmt.Sprintf("Transaction{From: %s, To: %s, Amount: %d, Data: %s}", tx.From.String(), tx.To.String(), tx.Amount, hex.EncodeToString(tx.Data))
}
