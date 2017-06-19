package playlist

import (
	"fmt"
	"io"
	"net/url"
)

const (
	playlistHeader = "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-TARGETDURATION:%d\n"
	playlistChunk  = "#EXTINF:%.3f,\n%s\n"
)

type Chunk struct {
	Length float64
	Path   string
}

type Writer struct {
	wr        io.Writer
	targetDur int
}

func NewWriter(wr io.Writer, targetDur int) *Writer {
	return &Writer{
		wr:        wr,
		targetDur: targetDur,
	}
}

func (w *Writer) WriteHeader() error {
	_, err := fmt.Fprintf(w.wr, playlistHeader, w.targetDur)
	return err
}

func (w *Writer) WriteChunk(c *Chunk) error {
	_, err := fmt.Fprintf(w.wr, playlistChunk, c.Length, c.Path)
	return err
}

func (w *Writer) WriteChunks(cs []*Chunk) error {
	for _, chunk := range cs {
		if err := w.WriteChunk(chunk); err != nil {
			return err
		}
	}
	return nil
}

func Path(root, streamId string, fragmentTag uint) string {
	streamId = url.QueryEscape(streamId)
	return fmt.Sprintf("%s/fragment.ts?stream=%s&fragment=%d", root, streamId, fragmentTag)
}
