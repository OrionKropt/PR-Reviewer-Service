package main

import (
	stdlog "log"
	"log/slog"
	"os"

	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/config"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/logger"
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
	_, log, err := initialze()
	if err != nil {
		stdlog.Fatal(err)
		return
	}
	log.Info("Initializing PR reviewer service")
}
