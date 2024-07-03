package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

// Peer represents a network peer.
type Peer struct {
	conn     net.Conn
	Outgoing bool
}

// NewPeer creates a new Peer instance.
func NewPeer(conn net.Conn, outgoing bool) *Peer {
	return &Peer{
		conn:     conn,
		Outgoing: outgoing,
	}
}

// Send sends data to the peer.
func (p *Peer) Send(data []byte) error {
	_, err := p.conn.Write(data)
	return err
}

// ReadLoop continuously reads data from the peer connection.
func (p *Peer) ReadLoop(rpcCh chan<- RPC) {
	buf := make([]byte, 4096)
	for {
		n, err := p.conn.Read(buf)
		if err == io.EOF {
			fmt.Printf("Connection closed by peer: %s\n", p.conn.RemoteAddr().String())
			break
		}
		if err != nil {
			fmt.Printf("Error reading from peer %s: %s\n", p.conn.RemoteAddr().String(), err)
			break
		}

		msg := make([]byte, n)
		copy(msg, buf[:n])
		rpcCh <- RPC{
			From:    p.conn.RemoteAddr().String(),
			Payload: bytes.NewReader(msg),
		}
	}
	p.conn.Close()
}

// Close closes the peer connection.
func (p *Peer) Close() error {
	return p.conn.Close()
}
