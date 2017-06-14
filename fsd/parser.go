package fsd

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"io"
)

// Parser holds the parser's internal state.
type Parser struct {
	rd *bufio.Reader
}

// NewParser returns a new Parser that reads from input rd.
func NewParser(rd io.Reader) *Parser {
	return &Parser{
		rd: bufio.NewReader(rd),
	}
}

// ParseFragment reads a single fragment from the stream and returns it.
func (p *Parser) ParseFragment() (*Fragment, error) {
	var fm Fragment

	// Length
	var rawLength [4]byte
	_, err := io.ReadFull(p.rd, rawLength[:])
	if err != nil {
		return nil, err
	}
	fm.Length = binary.LittleEndian.Uint32(rawLength[:])

	// Signature
	_, err = io.ReadFull(p.rd, fm.Signature[:])
	if err != nil {
		return nil, err
	}

	// Data
	fm.Data = make([]byte, int(fm.Length))
	_, err = io.ReadFull(p.rd, fm.Data)
	if err != nil {
		return nil, err
	}

	// HashedData
	fm.HashedData = sha256.Sum256(fm.Data)

	return &fm, nil
}
