package main

import (
	"database/sql"
	"fmt"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	logg := logger.New(cfg.LoggerLevel)

	if err := applyMigrations(logg, cfg.DatabaseURL); err != nil {
		logg.Error("Database migrations failed", "error", err)
		os.Exit(1)
	}

}

func applyMigrations(logg *slog.Logger, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB for migrations: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(os.DirFS("./migrations"))

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logg.Info("Database migrations applied successfully")
	return nil
}
