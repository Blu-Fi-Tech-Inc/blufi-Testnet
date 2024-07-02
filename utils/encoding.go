package utils

import (
	"encoding/gob"
	"io"
)

// Encoder is an interface for encoding a type to an io.Writer.
type Encoder interface {
	Encode(interface{}) error
}

// Decoder is an interface for decoding a type from an io.Reader.
type Decoder interface {
	Decode(interface{}) error
}

// GobTxEncoder encodes a Transaction to an io.Writer using GOB encoding.
type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w,
	}
}

func (enc *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(enc.w).Encode(tx)
}

// GobTxDecoder decodes a Transaction from an io.Reader using GOB decoding.
type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		r: r,
	}
}

func (dec *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(dec.r).Decode(tx)
}

// GobBlockEncoder encodes a Block to an io.Writer using GOB encoding.
type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

// GobBlockDecoder decodes a Block from an io.Reader using GOB decoding.
type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (dec *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(dec.r).Decode(b)
}
