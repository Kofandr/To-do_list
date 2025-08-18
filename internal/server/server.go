package server

import (
	"context"
	"github.com/Kofandr/To-do_list/config"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	addr string
	logg *slog.Logger
	db   *pgxpool.Pool
}

func New(logg *slog.Logger, cfg *config.Configuration, db *pgxpool.Pool) *Server {
	serverEcho := echo.New()

	return &Server{serverEcho, ":" + strconv.Itoa(cfg.Port), logg, db}
}

func (server *Server) Start() error {
	server.logg.Info("Starting server", "addr", server.addr)

	return server.echo.Start(server.addr)
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logg.Info("Shutting down server")

	return server.echo.Shutdown(ctx)
}
