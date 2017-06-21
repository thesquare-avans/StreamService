package transport

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

var (
	testPayload         = "this is a dummy payload, just for testing purposes."
	randomEncryptedHash = "3cf1fb36afaef5f05b737e74c7f25a1c7c997f522e2667bcaae8c590a8f83bbe762b8472aa5c918f586b98cdc62679bc2fd2399fc52d7eba7627ead0b5e47181eb66ec1f28b83a708be4939515ed8ba0b11ac33e7863215911f1a47be6c5f35f32310bbe2a7f5407b8d4801d62f251e80c2a53eaae6650211bef2e547bb124311f386dd836f11d54dbf23ee0899b0cb5dec7aa0033e3585c5fee5679b39c8e5777fea0000ce14d7bd3069e41a3f95a7e6e489be8df09dfde72416229aecbbb362aace0788f214e569e436340421e1d63ab14c5eff5ad0737d6009fd75d2b6f352adfd80af8cb4d9b932097eb142c560c21440cd59e6d8eda6d02932d9330fd6a"
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
