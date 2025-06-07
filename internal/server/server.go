package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/ayrtonvitor/http/internal/headers"
	"github.com/ayrtonvitor/http/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
	}

	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	headers := headers.GetDefaultHeaders(0)

	response.WriteStatusLine(conn, response.StatusCodeOk)
	response.WriteHeaders(conn, headers)
}
