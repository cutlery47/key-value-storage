package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultAddr            = "127.0.0.1:8080"
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	httpServ *http.Server

	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServ := &http.Server{
		Handler: handler,
		Addr:    defaultAddr,
	}

	server := &Server{
		httpServ: httpServ,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Run() {
	go s.Serve()

	// listening for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	s.Shutdown(sigChan)
}

func (s *Server) Serve() {
	if err := s.httpServ.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("http server error:", err)
	}
}

// graceful shutdown
func (s *Server) Shutdown(sigChan <-chan os.Signal) error {
	// waiting for shutdown signal to arrive
	<-sigChan

	log.Println("http server shutdown")

	// waiting for unfinished operations
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.httpServ.Shutdown(ctx)
}
