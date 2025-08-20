package handler

import (
	"fmt"
	"github.com/Kofandr/To-do_list/internal/repository"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Handler struct {
	db repository.Repository
}

func New(db repository.Repository) *Handler {
	return &Handler{
		db,
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
