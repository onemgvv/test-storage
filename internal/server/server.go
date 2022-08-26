package server

import (
	"context"
	"net/http"

	"test-storage/pkg/config"
)

type Server struct {
	httpServer *http.Server
	httpServe http.ServeMux
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.Timeouts.Read,
			WriteTimeout:   cfg.HTTP.Timeouts.Write,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
