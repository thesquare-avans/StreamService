package stream

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/thesquare-avans/StreamService/fsd"
)

type Server struct {
	l        net.Listener
	fragment int
	dir      string
	conn     net.Conn
	parser   *fsd.Parser
}

func NewServer(laddr, dir string) (*Server, error) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}
	return &Server{l: l, dir: dir}, nil
}

func (s *Server) WaitForClient() error {
	conn, err := s.l.Accept()
	if err != nil {
		return err
	}
	s.conn = conn
	s.parser = fsd.NewParser(s.conn)
	return err
}

func (s *Server) ReceiveSingle() (int, error) {
	fragment, err := s.parser.ParseFragment()
	if err != nil {
		return 0, err
	}

	file := fmt.Sprintf("StreamService-%d.mp4", s.fragment)
	s.fragment++
	fullPath := filepath.Join(s.dir, file)
	fd, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	_, err = fd.Write(fragment.Data)
	return len(fragment.Data), err
}

func (s *Server) Close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return s.l.Close()
}
