package handler

import (
	"context"
	"errors"
	"github.com/Kofandr/To-do_list/internal/appctx"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) DeleteUser(c echo.Context) error {
	return handler.HandleDelete(c, handler.db.DeleteUser, "User")
}

func (handler *Handler) DeleteTask(c echo.Context) error {
	return handler.HandleDelete(c, handler.db.DeleteTask, "Task")
}

func (handler *Handler) HandleDelete(c echo.Context, deleteFunc func(context.Context, int) error, entity string) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	id, err := parseIDParam(c)

	if err != nil {
		logg.Error("Invalid ID", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid ID"})
	}

	err = deleteFunc(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errResp := map[string]string{"err": "Not found"}

			logg.Error("Not found id", logger.ErrAttr(err))

			return c.JSON(http.StatusNotFound, errResp)
		}

		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": entity + " deleted"})
}
