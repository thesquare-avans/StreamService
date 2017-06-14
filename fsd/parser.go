package fsd

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
)

type Parser struct {
	rd *bufio.Reader
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{
		rd: bufio.NewReader(rd),
	}
}

func (p *Parser) ParseFragment() (*Fragment, error) {
	var fm Fragment

	// Length
	var rawLength [4]byte
	_, err := io.ReadFull(p.rd, rawLength[:])
	if err != nil {
		return nil, err
	}
	fm.Length = binary.LittleEndian.Uint32(rawLength[:])
	fmt.Printf("Raw bytes: %+v, length: %d\n", rawLength[:], fm.Length)

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
	fm.HashedData = sha256.Sum256(fm.Data)

	fmt.Printf("Length of data: %d\n", len(fm.Data))
	return &fm, nil
}
