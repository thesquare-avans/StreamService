package fsd

import (
	"crypto"
	"crypto/rsa"
	"strconv"

	"github.com/zhulik/go_mediainfo"
)

const (
	mediaInfoDuration = "Duration"
)

// Fragment represents a single signed fragment.
type Fragment struct {
	Length      uint32
	Signature   [256]byte
	Data        []byte
	HashedData  [32]byte
	Tag         uint
	VideoLength float64
}

// Verify verifies the fragment using its signature.
func (f *Fragment) Verify(publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, f.HashedData[:], f.Signature[:])
}

func (f *Fragment) GetVideoLength() (float64, error) {
	if f.VideoLength != 0 {
		return f.VideoLength, nil
	}
	info := mediainfo.NewMediaInfo()
	err := info.OpenMemory(f.Data)
	if err != nil {
		return 0, err
	}
	durationMs, err := strconv.ParseFloat(info.Get(mediaInfoDuration), 64)
	if err != nil {
		return 0, err
	}
	f.VideoLength = durationMs / 1000

	return f.VideoLength, nil
}
