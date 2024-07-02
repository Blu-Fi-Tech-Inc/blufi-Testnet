package network

import "github.com/blu-fi-tech-inc/boriqua_project/core"

// GetBlocksMessage represents a request to get blocks from a specific index to a maximum index (if To is 0).
type GetBlocksMessage struct {
	From uint32 // Starting index of blocks to fetch.
	To   uint32 // Ending index (inclusive) of blocks to fetch. If 0, fetch maximum available.
}

// BlocksMessage represents a response containing blocks.
type BlocksMessage struct {
	Blocks []*core.Block // List of blocks to send in response.
}

// GetStatusMessage represents a request to get status information.
type GetStatusMessage struct{}

// StatusMessage represents a response containing status information.
type StatusMessage struct {
	ID            string // ID of the server.
	Version       uint32 // Version of the server.
	CurrentHeight uint32 // Current height of the blockchain.
}
