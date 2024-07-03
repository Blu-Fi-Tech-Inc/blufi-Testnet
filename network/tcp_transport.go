package network

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

// TCPPeer represents a TCP peer.
type TCPPeer struct {
	conn     net.Conn
	Outgoing bool
}

// Send sends data to the peer.
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	if err != nil {
		logrus.WithError(err).Error("error sending data to peer")
	}
	return err
}

// readLoop continuously reads from the peer connection.
func (p *TCPPeer) readLoop(rpcCh chan RPC) {
	buf := make([]byte, 4096)
	for {
		n, err := p.conn.Read(buf)
		if err == io.EOF {
			continue // EOF is expected when connection closes
		}
		if err != nil {
			logrus.WithError(err).Error("read error from peer")
			continue
		}

		msg := make([]byte, n)
		copy(msg, buf[:n]) // Create a copy of the buffer to avoid concurrent access issues
		rpcCh <- RPC{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}
	}
}

// TCPTransport manages TCP connections and peers.
type TCPTransport struct {
	peerCh     chan *TCPPeer
	listenAddr string
	listener   net.Listener
}

// NewTCPTransport creates a new TCPTransport instance.
func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		peerCh:     peerCh,
		listenAddr: addr,
	}
}

// Start starts accepting incoming connections.
func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		logrus.WithError(err).Error("error starting TCP listener")
		return err
	}

	t.listener = ln

	go t.acceptLoop()

	logrus.Info("TCP transport started successfully")

	return nil
}

// acceptLoop continuously accepts incoming connections.
func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			logrus.WithError(err).Error("error accepting connection")
			continue
		}

		peer := &TCPPeer{
			conn: conn,
		}

		t.peerCh <- peer

		logrus.WithField("peer", conn.RemoteAddr()).Info("new peer connected")
	}
}
