package fairco

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
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
	_, err := p.rd.Read(rawLength[:])
	if err != nil {
		return nil, err
	}
	fm.Length = binary.LittleEndian.Uint32(rawLength[:])

	// Signature
	_, err = p.rd.Read(fm.Signature[:])
	if err != nil {
		return nil, err
	}

	// Data
	fm.Data = make([]byte, 0, int(fm.Length))
	_, err = io.ReadFull(p.rd, fm.Data)
	if err != nil {
		return nil, err
	}
	fm.HashedData = sha256.Sum256(fm.Data)
	return &fm, nil
}
