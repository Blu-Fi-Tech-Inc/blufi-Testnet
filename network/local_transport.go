package network

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

type LocalTransport struct {
	addr      net.Addr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[net.Addr]*LocalTransport
}

func NewLocalTransport(addr net.Addr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[net.Addr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	trans := tr.(*LocalTransport)
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = trans

	return nil
}

func (t *LocalTransport) SendMessage(to net.Addr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.addr == to {
		return nil
	}

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to unknown peer %s", t.addr, to)
	}

	select {
	case peer.consumeCh <- RPC{From: t.addr, Payload: bytes.NewReader(payload)}:
	default:
		return fmt.Errorf("failed to send message to %s", to)
	}

	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	var errs []error

	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to broadcast message: %v", errs)
	}

	return nil
}

func (t *LocalTransport) Addr() net.Addr {
	return t.addr
}
