package handler

import (
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/appctx"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/Kofandr/To-do_list/internal/repository/postgres"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	var task model.NewTask
	if err := c.Bind(&task); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(task); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	Exist, err := handler.db.UserExists(ctx, task.UserID)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	if !Exist {
		errResp := map[string]string{"err": "Not found User"}

		logg.Error("Not found User", logger.ErrAttr(err))

		return c.JSON(http.StatusNotFound, errResp)
	}

	id, err := handler.db.CreateTask(ctx, &task)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicate) {
			errResp := map[string]string{"err": fmt.Sprintf("Task with name '%s' already exists", task.Title)}

			logg.Warn("Duplicate task attempt", "name", task.Title, logger.ErrAttr(err))

			return c.JSON(http.StatusConflict, errResp)
		}

		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"ID task": id,
	})
}
