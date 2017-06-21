package stream

import (
	"net/http"
	"sync"
	"time"

	"github.com/thesquare-avans/StreamService/fsd"
)

type StreamServer struct {
	lock     sync.RWMutex
	current  *fsd.Fragment
	syncTime time.Time
}

func NewStreamServer() *StreamServer {
	return &StreamServer{}
}

func (s *StreamServer) PushFragment(f *fsd.Fragment) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.current = f
	s.syncTime = time.Now().Add(f.Duration)
}

func (s *StreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Send file
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(s.current.Data)
	if err != nil {
		panic(err)
	}
}
