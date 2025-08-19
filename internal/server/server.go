package server

import (
	"context"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/repository"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	addr string
	logg *slog.Logger
	db   repository.Repository
}

func New(logg *slog.Logger, cfg *config.Configuration, db repository.Repository) *Server {
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
