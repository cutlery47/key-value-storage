package server

import (
	"time"
)

// Option configuration pattern
type Option func(*Server)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.httpServ.Addr = addr
	}
}
