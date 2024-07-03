package api

import (
	"log"
	"net/http"

	"github.com/blu-fi-tech-inc/blufi-network/consensus"
	"github.com/blu-fi-tech-inc/blufi-network/core"
	"github.com/gorilla/mux"
)

// Server represents the API server.
type Server struct {
	api *API
}

// NewServer initializes a new API server instance.
func NewServer(txPool *core.TxPool, encoder core.Encoder, decoder core.Decoder, stakeManager *consensus.StakeManager, pos *consensus.PoS) *Server {
	api := NewAPI(txPool, encoder, decoder, stakeManager, pos)
	return &Server{api: api}
}

// Start initializes and starts the API server.
func (s *Server) Start(port string) {
	r := mux.NewRouter()
	s.api.RegisterRoutes(r)

	// Add logging middleware
	r.Use(LoggingMiddleware)

	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
