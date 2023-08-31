package main

import (
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/config"
	addToUserSegment "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/handlers/segment/addToUser"
	deleteSegment1 "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/handlers/segment/delete"
	saveSegment "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/handlers/segment/save"
	deleteUser "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/handlers/user/delete"
	saveUser "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/handlers/user/save"
	getUserSegments "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/handlers/user/segments"
	mwLogger "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/http-server/middleware/logger"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/logger/sl"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgresql.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	//log.Info("connect db", slog.String("env", cfg.Env))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/user", func(r chi.Router) {
		r.Post("/save", saveUser.New(log, storage))
		r.Delete("/delete", deleteUser.New(log, storage))
		r.Get("/segments", getUserSegments.New(log, storage))
	})

	router.Route("/segment", func(r chi.Router) {
		r.Post("/save", saveSegment.New(log, storage))
		r.Delete("/delete", deleteSegment1.New(log, storage))
		r.Post("/addToUser", addToUserSegment.New(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	//srv := &http.Server{
	//	Addr:         cfg.HTTPServer.Address,
	//	Handler:      router,
	//	ReadTimeout:  cfg.HTTPServer.Timeout,
	//	WriteTimeout: cfg.HTTPServer.Timeout,
	//	IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	//}

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
