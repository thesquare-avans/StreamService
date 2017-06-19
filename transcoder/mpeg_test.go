package transcoder

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const inputFile = "data/demo.mp4"

func TestMpegTsTranscoder(t *testing.T) {
	inputPath, err := filepath.Abs(inputFile)
	if err != nil {
		t.Error(err)
	}
	inputFd, err := os.Open(inputPath)
	if err != nil {
		t.Error(err)
	}
	defer inputFd.Close()

	tc, err := NewMpegTs(inputFd, ioutil.Discard)
	if err != nil {
		t.Error(err)
	}
	defer tc.Dispose()
	err = tc.Transcode()
	if err != nil {
		t.Error(err)
	}
}

func TestMpegTsTranscoderFile(t *testing.T) {
	inputPath, err := filepath.Abs(inputFile)
	if err != nil {
		t.Error(err)
	}

	tc := NewMpegTsFile(inputPath, ioutil.Discard)
	defer tc.Dispose()
	err = tc.Transcode()
	if err != nil {
		t.Error(err)
	}
}
