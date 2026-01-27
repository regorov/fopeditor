package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	httpapi "github.com/regorov/fopeditor/backend/internal/http"
	"github.com/regorov/fopeditor/backend/internal/render"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	port := getenv("PORT", "8080")
	var renderer render.Renderer
	if endpoint := os.Getenv("FOP_ENDPOINT"); endpoint != "" {
		logger.Info().Msgf("using FOP endpoint %s", endpoint)
		renderer = render.NewFOPRenderer(endpoint)
	} else {
		logger.Warn().Msg("FOP_ENDPOINT not set, falling back to stub renderer")
		renderer = render.NewStubRenderer()
	}
	api := httpapi.NewServer(renderer)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:          api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info().Msgf("HTTP server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server exited")
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	} else {
		logger.Info().Msg("server stopped")
	}
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
