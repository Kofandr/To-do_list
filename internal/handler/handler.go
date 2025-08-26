package handler

import (
	"fmt"
	"github.com/Kofandr/To-do_list/internal/repository"
	"github.com/Kofandr/To-do_list/internal/service/auth"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Handler struct {
	db      repository.Repository
	service *auth.Service
}

func New(db repository.Repository, service *auth.Service) *Handler {
	return &Handler{
		db,
		service,
	}
}

func parseIDParam(c echo.Context) (int, error) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return id, nil
}
