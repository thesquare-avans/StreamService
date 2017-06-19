package fsd

import (
	"crypto"
	"crypto/rsa"
)

// Fragment represents a single signed fragment.
type Fragment struct {
	Length     uint32
	Signature  [256]byte
	Data       []byte
	HashedData [32]byte
	Tag        uint
}

// Verify verifies the fragment using its signature.
func (f *Fragment) Verify(publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, f.HashedData[:], f.Signature[:])
}
