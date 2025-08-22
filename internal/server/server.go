package server

import (
	"context"
	"errors"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/appvalidator"
	"github.com/Kofandr/To-do_list/internal/handler"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/Kofandr/To-do_list/internal/middleware"
	"github.com/Kofandr/To-do_list/internal/repository"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
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

	serverEcho.Use(middleware.RequestLogger(logg))

	serverEcho.Validator = &appvalidator.CustomValidator{Validator: validator.New()}

	handler := handler.New(db)

	serverEcho.GET("/users", handler.GetUsers)
	serverEcho.POST("/users", handler.CreateUser)

	return &Server{serverEcho, ":" + strconv.Itoa(cfg.Port), logg, db}
}

func (server *Server) Start() error {
	server.logg.Info("Starting server", "addr", server.addr)

	err := server.echo.Start(server.addr)
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		server.logg.Info("Server stopped normally")
	} else {
		server.logg.Error("Server stopped with error", logger.ErrAttr(err))
	}

	return err
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logg.Info("Shutting down server")

	return server.echo.Shutdown(ctx)
}
