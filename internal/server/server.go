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
	"github.com/Kofandr/To-do_list/internal/service/auth"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo    *echo.Echo
	addr    string
	logg    *slog.Logger
	db      repository.Repository
	service *auth.Service
}

func New(logg *slog.Logger, cfg *config.Configuration, db repository.Repository, service *auth.Service) *Server {
	serverEcho := echo.New()

	serverEcho.Use(middleware.RequestLogger(logg))

	serverEcho.Validator = &appvalidator.CustomValidator{Validator: validator.New()}

	handler := handler.New(db, service)

	serverEcho.POST("/register", handler.RegisterUsers)
	serverEcho.POST("/login", handler.LoginUsers)
	serverEcho.POST("/refresh", handler.RefreshUsers)

	protected := serverEcho.Group("")
	protected.Use(middleware.JWTAuth(cfg))

	protected.POST("/tasks", handler.CreateTask)
	protected.GET("/user/tasks/:id", handler.GetTasksUser)
	protected.POST("/tasks/:id", handler.UpdateTask)
	protected.GET("/tasks/:id", handler.CompleteTask)
	protected.DELETE("/tasks/:id", handler.DeleteTask)

	return &Server{serverEcho, ":" + strconv.Itoa(cfg.Port), logg, db, service}
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
