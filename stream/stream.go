package stream

import (
	"fmt"
	"net"

	"github.com/thesquare-avans/StreamService/fsd"
)

type Server struct {
	l    *net.TCPListener
	pipe chan *fsd.Fragment
}

func NewServer(host string, port int, pipe chan *fsd.Fragment) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{l: l, pipe: pipe}, nil
}

func (s *Server) Run() error {
	client, err := s.l.AcceptTCP()
	if err != nil {
		return err
	}
	defer client.Close()
	// TODO: eventually close every other connection

	parser := fsd.NewParser(client)
	for {
		frag, err := parser.ParseFragment()
		if err != nil {
			return err
		}
		s.pipe <- frag
	}
}

func (s *Server) Close() error {
	close(s.pipe)
	return s.l.Close()
}
