package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/Kofandr/To-do_list/internal/repository/postgres"
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
	var exitCod int
	defer func() {
		if exitCod != 0 {
			os.Exit(exitCod)
		}
	}()

	cfg := config.MustLoad()
	logg := logger.New(cfg.LoggerLevel)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	errCh := make(chan error, 2)

	if err := applyMigrations(logg, cfg.DatabaseURL); err != nil {
		logg.Error("Database migrations failed", logger.ErrAttr(err))
		errCh <- err
	}

	db := postgres.New(pool)

	mainServer := server.New(logg, cfg, db)

	go func() {
		if err := mainServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error("Server crash")
			errCh <- err
		}
	}()

	var startErr error

	select {
	case sig := <-signalChan:
		logg.Info("Shutdown Request", "signal", sig.String())
	case err := <-errCh:
		startErr = err
		logg.Error("Server error", logger.ErrAttr(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShuttingDowntime)*time.Second)
	defer cancel()

	logg.Info("Shutting down...")

	if err := mainServer.Shutdown(ctx); err != nil {
		logg.Error("Shutdown failed", logger.ErrAttr(err))
	} else {
		logg.Info("Server stopped")
	}

	if startErr != nil {
		exitCod = 1
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
