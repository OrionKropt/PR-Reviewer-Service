package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/config"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/core"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/middleware"
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

	loggingHandler := middleware.LoggingMiddleware(log)
	rootMux := http.NewServeMux()

	teamMux := http.NewServeMux()
	teamMux.Handle("POST /add", loggingHandler(s.handleTeamAdd()))
	teamMux.Handle("GET /get", loggingHandler(s.handleTeamGet()))

	usersMux := http.NewServeMux()
	usersMux.Handle("POST /setIsActive", loggingHandler(s.handleUsersSetIsActive()))
	usersMux.Handle("GET /getReview", loggingHandler(s.handleUsersGetReview()))

	pullRequestMux := http.NewServeMux()
	pullRequestMux.Handle("POST /create", loggingHandler(s.handlePRCreate()))
	pullRequestMux.Handle("POST /merge", loggingHandler(s.handlePRMerge()))
	pullRequestMux.Handle("POST /reassign", loggingHandler(s.handlePRReassign()))

	rootMux.Handle("/team/", http.StripPrefix("/team", teamMux))
	rootMux.Handle("/users/", http.StripPrefix("/users", usersMux))
	rootMux.Handle("/pullRequest/", http.StripPrefix("/pullRequest", pullRequestMux))

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
