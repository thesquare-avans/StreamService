package transport

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

var (
	testPayload         = "this is a dummy payload, just for testing purposes."
	randomEncryptedHash = "PPH7Nq+u9fBbc350x/JaHHyZf1IuJme8qujFkKj4O752K4RyqlyRj1hrmM3GJnm8L9I5n8Utfrp2J+rQteRxgetm7B8ouDpwi+STlRXti6CxGsM+eGMhWRHxpHvmxfNfMjELvip/VAe41IAdYvJR6AwqU+quZlAhG+8uVHuxJDEfOG3YNvEdVNvyPuCJmwy13seqADPjWFxf7lZ5s5yOV3f+oAAM4U170waeQaP5Wn5uSJvo3wnf3nJBYimuy7s2KqzgeI8hTlaeQ2NAQh4dY6sUxe/1rQc31gCf110rbzUq39gK+MtNm5Mgl+sULFYMIUQM1Z5tjtptApMtkzD9ag=="
)

func TestPayloadVerifySignature(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}

	payload := Payload{Payload: testPayload}
	err = payload.Sign(privateKey)
	if err != nil {
		t.Error(err)
	}
	err = payload.Verify(&privateKey.PublicKey)
	if err != nil {
		t.Error(err)
	}

	// Poison signature and try to verify again
	payload.Signature = randomEncryptedHash
	err = payload.Verify(&privateKey.PublicKey)
	// err should be ErrVerification
	if err != rsa.ErrVerification {
		t.Error(err)
	}
}
