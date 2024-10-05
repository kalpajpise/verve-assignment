package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler, host, port string, timeout time.Duration) (*Server, error) {
	switch {
	case port == "":
		return nil, fmt.Errorf("empty server port")
	case handler == nil:
		return nil, fmt.Errorf("empty handler")
	case timeout < 1:
		return nil, fmt.Errorf("invalid timeout %d", timeout)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	if err := http2.ConfigureServer(server, nil); err != nil {
		log.Printf("can't configure http server: %s", err)
	}

	return &Server{
		server: server,
	}, nil
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
}
