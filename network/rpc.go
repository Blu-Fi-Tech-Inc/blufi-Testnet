package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3
	MessageTypeStatus    MessageType = 0x4
	MessageTypeGetStatus MessageType = 0x5
	MessageTypeBlocks    MessageType = 0x6
)

// RPC represents a Remote Procedure Call.
type RPC struct {
	From    net.Addr // Address of the sender.
	Payload io.Reader // Payload of the RPC.
}

// Message represents a network message.
type Message struct {
	Header MessageType // Type of message.
	Data   []byte      // Data payload.
}

// NewMessage creates a new Message instance.
func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

// Bytes serializes the Message into bytes.
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

// DecodedMessage represents a decoded RPC message.
type DecodedMessage struct {
	From net.Addr // Address of the sender.
	Data interface{} // Decoded data payload.
}

// RPCDecodeFunc defines a function type for decoding RPC messages.
type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

// DefaultRPCDecodeFunc decodes an RPC message based on its type.
func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, fmt.Errorf("failed to decode transaction message from %s: %s", rpc.From, err)
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil

	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, fmt.Errorf("failed to decode block message from %s: %s", rpc.From, err)
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil

	case MessageTypeGetStatus:
		return &DecodedMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil

	case MessageTypeStatus:
		statusMessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
			return nil, fmt.Errorf("failed to decode status message from %s: %s", rpc.From, err)
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil

	case MessageTypeGetBlocks:
		getBlocks := new(GetBlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocks); err != nil {
			return nil, fmt.Errorf("failed to decode get blocks message from %s: %s", rpc.From, err)
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: getBlocks,
		}, nil

	case MessageTypeBlocks:
		blocks := new(BlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blocks); err != nil {
			return nil, fmt.Errorf("failed to decode blocks message from %s: %s", rpc.From, err)
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blocks,
		}, nil

	default:
		return nil, fmt.Errorf("invalid message header %x from %s", msg.Header, rpc.From)
	}
}

// RPCProcessor defines an interface for processing decoded RPC messages.
type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

func init() {
	gob.Register(core.Transaction{})
	gob.Register(core.Block{})
	gob.Register(&GetStatusMessage{})
	gob.Register(&StatusMessage{})
	gob.Register(&GetBlocksMessage{})
	gob.Register(&BlocksMessage{})
}
