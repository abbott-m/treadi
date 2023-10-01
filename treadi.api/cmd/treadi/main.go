package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"treadi.api/internal/config"
	"treadi.api/internal/handlers"
	"treadi.api/internal/hub"
)

func main() {
	cg := config.New()
	flag.StringVar(&cg.Port, "port", "42069", "override port from env")
	flag.Parse()

	lh := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(lh)

	if err := run(cg, logger); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(config *config.Config, logger *slog.Logger) error {
	h := hub.NewHub()
	go h.Run()
	hh := handlers.NewHubHandler(h, logger)

	http.HandleFunc("/ws", hh.Serve)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.Port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Debug("listening on: ", slog.String("port", config.Port))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to listen: ", err)
		return err
	}

	return nil
}
