package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/config"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/core"
)

type Server struct {
	httpServer   *http.Server
	prRevService *core.PRReviewerService
	log          *slog.Logger
}

func NewServer(cfg *config.Config, log *slog.Logger) *Server {
	s := &Server{
		prRevService: core.NewPRReviewerService(),
		log:          log,
	}

	rootMux := http.NewServeMux()
	
	s.httpServer = &http.Server{
		Addr:    cfg.BindAddr,
		Handler: rootMux,
	}

	return s
}

func (s *Server) Start() error {
	s.log.Info("HTTP server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("HTTP server stopping")
	return s.httpServer.Shutdown(ctx)
}
