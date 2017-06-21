package fsd

import (
	"crypto"
	"crypto/rsa"
	"strconv"
	"time"

	"github.com/zhulik/go_mediainfo"
)

const (
	mediaInfoDuration = "Duration"
)

// Fragment represents a single signed fragment.
type Fragment struct {
	Length     uint32
	Signature  [256]byte
	Data       []byte
	HashedData [32]byte
	Duration   time.Duration
}

// Verify verifies the fragment using its signature.
func (f *Fragment) Verify(publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, f.HashedData[:], f.Signature[:])
}

func (f *Fragment) GetDuration() (time.Duration, error) {
	if f.Duration != 0 {
		return f.Duration, nil
	}
	info := mediainfo.NewMediaInfo()
	err := info.OpenMemory(f.Data)
	if err != nil {
		return 0, err
	}
	duration, err := strconv.Atoi(info.Get(mediaInfoDuration))
	if err != nil {
		return 0, err
	}
	f.Duration = time.Duration(duration) * time.Millisecond

	return f.Duration, nil
}
