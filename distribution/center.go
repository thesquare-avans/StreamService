package distribution

import (
	"container/list"
	"errors"
	"sync"

	"github.com/thesquare-avans/StreamService/fsd"
)

const MaxFragmentsInBuffer = 5

var (
	ErrStreamExists     = errors.New("distribution: stream already exists")
	ErrStreamNotExists  = errors.New("distribution: stream doesn't exists")
	ErrFragmentNotFound = errors.New("distribution: fragment not found")
)

type fragmentBuffer struct {
	lock     sync.RWMutex
	buffer   *list.List
	lastTag  uint
	mediaSeq int
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
	if _, exists := c.streams[id]; !exists {
		return ErrStreamNotExists
	}
	// TODO: maybe lock the stream lock.
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
		s.buffer.Remove(s.buffer.Front())
		s.mediaSeq++
	}
	// NOTE: this alters the tag of Fragment f.
	f.Tag = s.lastTag
	s.lastTag++
	s.buffer.PushBack(f)
	return nil
}

func (c *Center) GetFragmentFromStream(id string, tag uint) (*fsd.Fragment, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	s, exists := c.streams[id]
	if !exists {
		return nil, ErrStreamNotExists
	}
	s.lock.RLock()
	defer s.lock.RUnlock()

	for e := s.buffer.Front(); e != nil; e = e.Next() {
		frag, _ := e.Value.(*fsd.Fragment)
		if frag.Tag == tag {
			return frag, nil
		}
	}

	return nil, ErrFragmentNotFound
}

func (c *Center) GetFragmentsFromStream(id string, max int) ([]*fsd.Fragment, int, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	s, exists := c.streams[id]
	if !exists {
		return nil, 0, ErrStreamNotExists
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	bufLen := s.buffer.Len()
	if max >= 0 && bufLen > max {
		bufLen = max
	}
	buf := make([]*fsd.Fragment, bufLen)

	e := s.buffer.Back()
	for i := 1; i < bufLen; i++ {
		e = e.Prev()
	}
	var i int
	for ; e != nil; e = e.Next() {
		frag, _ := e.Value.(*fsd.Fragment)
		buf[i] = frag
		i++
	}

	return buf, s.mediaSeq, nil
}
