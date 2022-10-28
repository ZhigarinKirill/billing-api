package apiserver

import (
	"context"
	"net/http"
	"time"
)

// Server - structure for server app
type HTTPServer struct {
	httpServer *http.Server
}

// NewServer - server constructor
func NewServer(port string, handler http.Handler) *HTTPServer {
	return &HTTPServer{httpServer: &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}}
}

// Server starting
func (s *HTTPServer) Start() error {
	return s.httpServer.ListenAndServe()
}

// Server shutdown
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
