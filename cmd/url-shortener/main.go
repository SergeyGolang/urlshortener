package main

import (
	"log/slog"
	"net/http"
	"os"
	"urlshortener/internal/config"
	redirect "urlshortener/internal/http-server/handlers/redirect"
	delete "urlshortener/internal/http-server/handlers/url/delete"
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

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		// можно return, но при os.Exit(1)-видно, что приложение упало с ошибкой.
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, chi render
	router := chi.NewRouter()

	// TODO: middleware
	// присваивание каждому запросу своего ID
	router.Use(middleware.RequestID)

	// логирование запросов в хендлере
	router.Use(mwLogger.New(log))

	// восстановление паники в хендлере(из за одного панического запроса не должно падать всё приложение)
	router.Use(middleware.Recoverer)

	// подключение красивых URL к роутеру
	router.Use(middleware.URLFormat)

	// общий префикс для группы модифицирующих операция url и подключение авторизации
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))
	// TODO: run server: main

	log.Info("starting server", slog.String("address", cfg.Addres))

	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout, //время чтения запроса
		WriteTimeout: cfg.HTTPServer.Timeout, // время отправки ответа и прочтения клиентом
		IdleTimeout:  cfg.HTTPServer.Idletimeout,
	}

	// это блокирующий вызов и код после не должен быть выполнен, если выполнен, то произошла ошибка и сервер остановился
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
