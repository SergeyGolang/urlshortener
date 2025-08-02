package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	ssogrpc "urlshortener/internal/clients/auth/grpc"
	"urlshortener/internal/config"
	delete "urlshortener/internal/http-server/handlers/url/delete"
	redirect "urlshortener/internal/http-server/handlers/url/redirect"
	save "urlshortener/internal/http-server/handlers/url/save"
	"urlshortener/internal/storage/sqlite"

	mwLogger "urlshortener/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("starting urlshortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", slog.Any("error", err))
	}

	ssoClient.IsAdmin(context.Background(), 1)

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, chi render
	router := chi.NewRouter()

	// TODO: add middleware
	// Assigns a unique ID to each request for tracking
	router.Use(middleware.RequestID)

	// Logs all incoming requests (handler-level logging)
	router.Use(mwLogger.New(log))

	// Recovers from panics to prevent app-wide crashes
	router.Use(middleware.Recoverer)

	// Enables clean URL routing (e.g., /resource/{id})
	router.Use(middleware.URLFormat)

	router.Get("/{alias}", redirect.New(log, storage))
	// Enables BasicAuth
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	// TODO: run server: main

	log.Info("starting server", slog.String("address", cfg.Addres))

	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Idletimeout,
	}

	// This is a blocking call - code below will only execute if server fails
	// If you see this log, it means the server crashed unexpectedly
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
