package handler

import (
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/appctx"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/repository/postgres"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	var user model.NewUser
	if err := c.Bind(&user); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", "err", err)

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(user); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", "err", err)

		return c.JSON(http.StatusBadRequest, errResp)
	}

	id, err := handler.db.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicate) {
			errResp := map[string]string{"err": fmt.Sprintf("User with name '%s' already exists", user.Name)}

			logg.Warn("Duplicate user attempt", "name", user.Name, "err", err)

			return c.JSON(http.StatusConflict, errResp) // 409 Conflict
		}

		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", "err", err)

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"ID category": id,
	})
}
