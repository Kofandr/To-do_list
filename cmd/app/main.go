package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/Kofandr/To-do_list/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	logg := logger.New(cfg.LoggerLevel)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := applyMigrations(logg, cfg.DatabaseURL); err != nil {
		logg.Error("Database migrations failed", "error", err)
		os.Exit(1)
	}
	mainServer := server.New(logg, cfg, pool)

	go func() {
		if err := mainServer.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server crash")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShuttingDowntime)*time.Second)
	defer cancel()

	logg.Info("Shutting down...")

	if err := mainServer.Shutdown(ctx); err != nil {
		logg.Error("Shutdown failed", "error", err)
	} else {
		logg.Info("Server stopped")
	}

}

func applyMigrations(logg *slog.Logger, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB for migrations: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logg.Info("Database migrations applied successfully")
	return nil
}
