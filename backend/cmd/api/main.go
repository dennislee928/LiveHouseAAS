package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/api"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/config"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/cache"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/db"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/notification"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	pg, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	r := cache.NewRedis(cfg.RedisURL)
	defer r.Close()

	var notifSvc notification.Service
	if cfg.SMTPHost != "" {
		notifCfg := notification.Config{
			SMTPHost:     cfg.SMTPHost,
			SMTPPort:     cfg.SMTPPort,
			SMTPUser:     cfg.SMTPUser,
			SMTPPassword: cfg.SMTPPassword,
			FromEmail:    cfg.FromEmail,
		}
		notifSvc = notification.NewSMTPService(notifCfg)
		slog.Info("email service configured", "host", cfg.SMTPHost, "port", cfg.SMTPPort)
	} else {
		notifSvc = notification.NewConsoleService()
		slog.Info("email service using console (no SMTP configured)")
	}

	router := api.NewRouter(cfg, pg, r, notifSvc)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("API server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-quit
	slog.Info("shutting down server", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exited")
}
