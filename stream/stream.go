package stream

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/thesquare-avans/StreamService/distribution"
)

var (
	ErrNotFound = errors.New("bad request or resource is unavailable")
)

type jsonMediaSeq struct {
	Ok       bool `json:"ok"`
	MediaSeq uint `json:"mediaSequence"`
}

type jsonError struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type StreamServer struct {
	center *distribution.Center
}

func NewStreamServer(center *distribution.Center) *StreamServer {
	return &StreamServer{
		center: center,
	}
}

func (s *StreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if r.URL.Path == "/stream/mediaSequence" && r.Method == http.MethodGet {
		s.ServeMediaSeq(w, r)
	} else if r.URL.Path == "/stream/fragment.mp4" && r.Method == http.MethodGet {
		s.ServeFragment(w, r)
	} else {
		serveJsonError(w, http.StatusNotFound, ErrNotFound)
	}
}

func (s *StreamServer) ServeMediaSeq(w http.ResponseWriter, r *http.Request) {
	streamId := r.URL.Query().Get("stream")

	mediaSeq, err := s.center.GetMediaSeqFromStream(streamId)
	if err != nil {
		var status int
		if err == distribution.ErrStreamNotExists || err == distribution.ErrNoFragments {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		serveJsonError(w, status, err)
		return
	}

	mediaSeqStr, _ := json.Marshal(jsonMediaSeq{Ok: true, MediaSeq: mediaSeq})
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(mediaSeqStr)))
	w.WriteHeader(http.StatusOK)
	w.Write(mediaSeqStr)
}

func (s *StreamServer) ServeFragment(w http.ResponseWriter, r *http.Request) {
	streamId := r.URL.Query().Get("stream")
	mediaSeqParam := r.URL.Query().Get("mediaSequence")
	mediaSeq, err := strconv.Atoi(mediaSeqParam)
	if err != nil {
		serveJsonError(w, http.StatusBadRequest, err)
		return
	}

	fragment, err := s.center.GetFragmentFromStream(streamId, uint(mediaSeq))
	if err != nil {
		var status int
		if err == distribution.ErrInvalidSeq || err == distribution.ErrNoFragments || err == distribution.ErrNotFound {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		serveJsonError(w, status, err)
		return
	}

	w.Header().Add("Content-Type", "video/mp4")
	w.Header().Add("Content-Length", strconv.Itoa(len(fragment.Data)))
	w.WriteHeader(http.StatusOK)
	w.Write(fragment.Data)
}

func serveJsonError(w http.ResponseWriter, status int, err error) {
	errStr, _ := json.Marshal(jsonError{Error: err.Error()})
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(errStr)))
	w.WriteHeader(status)
	w.Write(errStr)
}
