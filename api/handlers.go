package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blu-fi-tech-inc/blufi-network/consensus"
	"github.com/blu-fi-tech-inc/blufi-network/core"
	"github.com/blu-fi-tech-inc/blufi-network/types"
	"github.com/blu-fi-tech-inc/blufi-network/utils"
	"github.com/gorilla/mux"
)

// API struct holds the necessary dependencies for API handlers.
type API struct {
	txPool       *core.TxPool
	encoder      core.Encoder
	decoder      core.Decoder
	stakeManager *consensus.StakeManager
	pos          *consensus.PoS
}

// NewAPI initializes a new API instance.
func NewAPI(txPool *core.TxPool, encoder core.Encoder, decoder core.Decoder, stakeManager *consensus.StakeManager, pos *consensus.PoS) *API {
	return &API{
		txPool:       txPool,
		encoder:      encoder,
		decoder:      decoder,
		stakeManager: stakeManager,
		pos:          pos,
	}
}

// RegisterRoutes registers all API routes with the provided router.
func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/transactions", a.handleNewTransaction).Methods("POST")
	r.HandleFunc("/transactions/pending", a.handleGetPendingTransactions).Methods("GET")
	r.HandleFunc("/blocks", a.handleNewBlock).Methods("POST")
	r.HandleFunc("/stake", a.handleAddStake).Methods("POST")
	r.HandleFunc("/stake/{address}", a.handleGetStake).Methods("GET")
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
	var block types.Block
	err := json.NewDecoder(r.Body).Decode(&block)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding block: %v", err), http.StatusBadRequest)
		return
	}

	// Validate and add the block using PoS
	if !a.pos.AddBlock(&block) {
		http.Error(w, "block validation failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(block)
}

// handleAddStake handles incoming POST requests to add stake to an address.
func (a *API) handleAddStake(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Address string `json:"address"`
		Amount  uint64 `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding payload: %v", err), http.StatusBadRequest)
		return
	}

	err = a.stakeManager.AddStake(payload.Address, payload.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("error adding stake: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleGetStake handles incoming GET requests to fetch the stake of an address.
func (a *API) handleGetStake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	stake, err := a.stakeManager.GetStake(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching stake: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stake)
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
