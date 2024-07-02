package consensus

import (
	"fmt"
	"sync"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
)

// Consensus represents the consensus mechanism.
type Consensus struct {
	mu       sync.Mutex
	chain    []*core.Block
	peers    []*network.Peer
	proposal *core.Block // The current proposed block for consensus
}

// NewConsensus creates a new Consensus instance.
func NewConsensus() *Consensus {
	return &Consensus{
		chain: []*core.Block{},
		peers: []*network.Peer{},
	}
}

// AddPeer adds a new peer to the consensus network.
func (c *Consensus) AddPeer(peer *network.Peer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.peers = append(c.peers, peer)
}

// ProposeBlock proposes a new block for consensus.
func (c *Consensus) ProposeBlock(block *core.Block) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.proposal = block
}

// FinalizeBlock finalizes the proposed block and adds it to the chain.
func (c *Consensus) FinalizeBlock() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.proposal == nil {
		return fmt.Errorf("no block proposal to finalize")
	}

	c.chain = append(c.chain, c.proposal)
	c.proposal = nil // Clear the current proposal after finalization

	return nil
}

// GetChain returns the current blockchain.
func (c *Consensus) GetChain() []*core.Block {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.chain
}

// HandleRPC handles incoming RPC messages related to consensus.
func (c *Consensus) HandleRPC(rpc network.RPC) error {
	// Example of handling different types of RPC messages related to consensus
	switch rpc.Type {
	case network.RPCTypes.Proposal:
		// Process proposal from peer
		// For example:
		// c.ProposeBlock(rpc.Payload)
		return nil
	case network.RPCTypes.Vote:
		// Process vote from peer
		// For example:
		// c.ProcessVote(rpc.Payload)
		return nil
	default:
		return fmt.Errorf("unknown RPC type received: %s", rpc.Type)
	}
}

// Start starts the consensus process.
func (c *Consensus) Start() {
	// Example of starting the consensus process
	// For example, listen for RPC messages, handle proposals, votes, etc.
	// You can implement your specific logic here based on your consensus algorithm.
}

// Stop stops the consensus process.
func (c *Consensus) Stop() {
	// Example of stopping the consensus process
	// For example, stop listening for RPC messages, clean up resources, etc.
}

// Example function to demonstrate usage
func ExampleUsage() {
	consensus := NewConsensus()

	// Example of adding peers
	peer1 := &network.Peer{} // Initialize your peer here
	peer2 := &network.Peer{} // Initialize another peer
	consensus.AddPeer(peer1)
	consensus.AddPeer(peer2)

	// Example of proposing and finalizing blocks
	block := &core.Block{} // Initialize your block here
	consensus.ProposeBlock(block)
	consensus.FinalizeBlock()

	// Example of retrieving the blockchain
	chain := consensus.GetChain()
	fmt.Printf("Current blockchain length: %d\n", len(chain))
}
