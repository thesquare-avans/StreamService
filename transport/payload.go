package transport

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

// Payload contains the payload itself and its signature.
type Payload struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

// Sign signs the payload.
func (p *Payload) Sign(privateKey *rsa.PrivateKey) error {
	hashed := sha256.Sum256([]byte(p.Payload))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return err
	}
	p.Signature = hex.EncodeToString(signature)
	return nil
}

// Verify verifies the payload using its signature.
func (p *Payload) Verify(publicKey *rsa.PublicKey) error {
	signature, err := hex.DecodeString(p.Signature)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(p.Payload))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
}
