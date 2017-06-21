package distribution

import (
	"container/list"
	"errors"
	"log"
	"sync"

	"github.com/thesquare-avans/StreamService/fsd"
)

const MaxFragmentsInBuffer = 5

var (
	ErrStreamExists    = errors.New("distribution: stream already exists")
	ErrStreamNotExists = errors.New("distribution: stream doesn't exists")
	ErrNotFound        = errors.New("distribution: fragment not found")
	ErrNoFragments     = errors.New("distribution: no fragments in stream")
	ErrInvalidSeq      = errors.New("distribution: invalid media sequence")
)

type fragmentBuffer struct {
	lock     sync.RWMutex
	buffer   *list.List
	mediaSeq uint
}

type Center struct {
	lock    sync.RWMutex
	streams map[string]*fragmentBuffer
}

func NewCenter() *Center {
	return &Center{
		streams: make(map[string]*fragmentBuffer),
	}
}

func (c *Center) NewStream(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.streams[id]; exists {
		return ErrStreamExists
	}
	c.streams[id] = &fragmentBuffer{buffer: list.New()}
	return nil
}

func (c *Center) DeleteStream(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	s, exists := c.streams[id]
	if !exists {
		return ErrStreamNotExists
	}
	s.lock.Lock()
	delete(c.streams, id)
	return nil
}

func (c *Center) PushToStream(id string, f *fsd.Fragment) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	s, exists := c.streams[id]
	if !exists {
		return ErrStreamNotExists
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.buffer.Len()+1 > MaxFragmentsInBuffer {
		s.buffer.Remove(s.buffer.Back())
	}
	s.buffer.PushFront(f)
	s.mediaSeq++
	log.Println("Pushed frame to", id, "media sequence", s.mediaSeq, "in buffer", s.buffer.Len())
	return nil
}

func (c *Center) GetMediaSeqFromStream(id string) (uint, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	s, exists := c.streams[id]
	if !exists {
		return 0, ErrStreamNotExists
	}
	s.lock.RLock()
	defer s.lock.RUnlock()

	mediaSeq := s.mediaSeq
	if mediaSeq < 1 {
		return 0, ErrNoFragments
	}

	return mediaSeq, nil
}

func (c *Center) GetFragmentFromStream(id string, mediaSeq uint) (*fsd.Fragment, error) {
	if mediaSeq < 1 {
		return nil, ErrInvalidSeq
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	s, exists := c.streams[id]
	if !exists {
		return nil, ErrStreamNotExists
	}
	s.lock.RLock()
	defer s.lock.RUnlock()

	i := s.mediaSeq
	e := s.buffer.Front()
	for {
		if i == 0 || e == nil {
			break
		} else if i == mediaSeq {
			return e.Value.(*fsd.Fragment), nil
		}
		i--
		e = e.Next()
	}

	return nil, ErrNotFound
}
