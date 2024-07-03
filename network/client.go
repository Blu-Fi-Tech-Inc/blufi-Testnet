package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

// Client represents a network client.
type Client struct {
	conn net.Conn
}

// NewClient creates a new network client.
func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// SendData sends data to the connected server.
func (c *Client) SendData(data interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}

	if _, err := c.conn.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// ReceiveData receives data from the connected server.
func (c *Client) ReceiveData() (interface{}, error) {
	var data interface{}
	dec := gob.NewDecoder(c.conn)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// Close closes the client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
