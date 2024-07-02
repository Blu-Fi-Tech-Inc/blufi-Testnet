package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
	"github.com/blu-fi-tech-inc/boriqua_project/types"
	"github.com/blu-fi-tech-inc/boriqua_project/utils"
	"github.com/gorilla/mux"
)

// API struct holds the necessary dependencies for API handlers.
type API struct {
	txPool  *core.TxPool
	encoder core.Encoder
	decoder core.Decoder
}

// NewAPI initializes a new API instance.
func NewAPI(txPool *core.TxPool, encoder core.Encoder, decoder core.Decoder) *API {
	return &API{
		txPool:  txPool,
		encoder: encoder,
		decoder: decoder,
	}
}

// RegisterRoutes registers all API routes with the provided router.
func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/transactions", a.handleNewTransaction).Methods("POST")
	r.HandleFunc("/transactions/pending", a.handleGetPendingTransactions).Methods("GET")
	r.HandleFunc("/blocks", a.handleNewBlock).Methods("POST")
}

// handleNewTransaction handles incoming POST requests to create a new transaction.
func (a *API) handleNewTransaction(w http.ResponseWriter, r *http.Request) {
	var tx core.Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Add transaction to the transaction pool
	a.txPool.Add(&tx)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx.Hash(core.TxHasher{}))
}

// handleGetPendingTransactions handles incoming GET requests to fetch pending transactions.
func (a *API) handleGetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	pendingTxs := a.txPool.Pending()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pendingTxs)
}

// handleNewBlock handles incoming POST requests to create a new block.
func (a *API) handleNewBlock(w http.ResponseWriter, r *http.Request) {
	var block core.Block
	err := json.NewDecoder(r.Body).Decode(&block)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding block: %v", err), http.StatusBadRequest)
		return
	}

	// Example: Validate block, sign, add to blockchain, etc.
	// For demonstration purposes, this example only echoes back the received block.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(block)
}

// Middleware example: Logging middleware for request tracing.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		utils.Log(fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path))
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
