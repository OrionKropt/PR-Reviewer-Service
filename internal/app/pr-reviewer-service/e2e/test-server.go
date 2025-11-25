package e2e

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/config"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/logger"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/server"
)

func StartTestServer(t *testing.T) (srv *server.Server, ts *httptest.Server) {
	t.Helper()

	cfg := &config.Config{
		BindAddr:           "8080",
		ReadHandlerTimeout: 5 * time.Second,
		ReadTimeout:        5 * time.Second,
		WriteTimeout:       5 * time.Second,
		IdleTimeout:        10 * time.Second,
	}

	log := logger.NewTest()

	srv = server.NewServer(cfg, log)
	ts = httptest.NewServer(srv.HttpHandler())
	return
}
