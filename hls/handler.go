package hls

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/thesquare-avans/StreamService/distribution"
	"github.com/thesquare-avans/StreamService/playlist"
)

type Handler struct {
	center    *distribution.Center
	root      string
	targetDur int
}

func NewHandler(center *distribution.Center, root string, targetDuration int) *Handler {
	return &Handler{center: center, root: root, targetDur: targetDuration}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Path:", r.URL.Path)

	w.Header().Add("Access-Control-Allow-Origin", "*")

	if r.URL.Path == "/live.m3u8" {
		h.ServePlaylist(w, r)
		return
	}
	h.ServeFragment(w, r)
}

func (h *Handler) ServePlaylist(w http.ResponseWriter, r *http.Request) {
	streamId := r.URL.Query().Get("stream")
	fragments, mediaSeq, err := h.center.GetFragmentsFromStream(streamId, 5)
	if err != nil {
		if err == distribution.ErrStreamNotExists {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "error: %s\r\n", err)
		return
	}

	var buffer bytes.Buffer
	pl := playlist.NewWriter(&buffer, h.targetDur, mediaSeq)
	err = pl.WriteHeader()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s\r\n", err)
		return
	}

	var chunk playlist.Chunk
	for _, fragment := range fragments {
		chunk.Path = playlist.Path(h.root, streamId, fragment.Tag)
		chunk.Length = fragment.VideoLength
		err = pl.WriteChunk(&chunk)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %s\r\n", err)
			return
		}
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	w.WriteHeader(http.StatusOK)
	w.Write(buffer.Bytes())
}

func (h *Handler) ServeFragment(w http.ResponseWriter, r *http.Request) {
	streamId := r.URL.Query().Get("stream")
	rawTag := r.URL.Query().Get("fragment")
	fragmentTag, err := strconv.Atoi(rawTag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s\r\n", err)
		return
	}
	fragment, err := h.center.GetFragmentFromStream(streamId, uint(fragmentTag))
	if err != nil {
		if err == distribution.ErrStreamNotExists || err == distribution.ErrFragmentNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "error: %s\r\n", err)
		return
	}
	w.Header().Add("Content-Type", "video/MP2T")
	w.WriteHeader(http.StatusOK)
	w.Write(fragment.Data)
}
