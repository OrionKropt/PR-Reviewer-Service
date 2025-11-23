package main

import (
	stdlog "log"
	"log/slog"
	"os"

	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/config"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/logger"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/server"
)

func initialze() (cfg *config.Config, log *slog.Logger, err error) {
	cfg = config.NewConfig()
	err = cfg.ReadConfig()
	if err != nil {
		return nil, nil, err
	}
	logHandler := logger.NewLogHandler(os.Stdout, logger.LogHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: logger.ParseLevel(cfg.LogLevel),
		},
	})
	log = slog.New(logHandler)
	return cfg, log, nil
}

func main() {
	cfg, log, err := initialze()
	if err != nil {
		stdlog.Fatal(err)
		return
	}
	log.Info("Initializing PR reviewer service")

	serv := server.NewServer(cfg, log)

	if err = serv.Start(); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
