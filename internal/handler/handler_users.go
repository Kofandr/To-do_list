package handler

import (
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/appctx"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/Kofandr/To-do_list/internal/repository/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	var user model.NewUser
	if err := c.Bind(&user); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(user); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	id, err := handler.db.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicate) {
			errResp := map[string]string{"err": fmt.Sprintf("User with Username '%s' already exists", user.Username)}

			logg.Warn("Duplicate user attempt", "Username", user.Username, logger.ErrAttr(err))

			return c.JSON(http.StatusConflict, errResp)
		}

		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"ID user": id,
	})
}

func (handler *Handler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	users, err := handler.db.GetUsers(ctx)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, users)
}

func (handler *Handler) RegisterUsers(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	var user model.NewUser
	if err := c.Bind(&user); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(user); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	id, err := handler.service.Register(&user, ctx)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"ID user": id,
	})
}

func (handler *Handler) LoginUsers(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	var user model.NewUser
	if err := c.Bind(&user); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(user); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	result, err := handler.service.Login(&user, ctx, logg)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, result)
}

func (handler *Handler) RefreshUsers(c echo.Context) error {
	ctx := c.Request().Context()

	logg := appctx.LoggerFromContext(ctx)

	var rec model.Tokens
	if err := c.Bind(&rec); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(rec); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	tokens, err := handler.service.Refresh(&rec, ctx)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, tokens)
}
func (handler *Handler) BindTelegram(c echo.Context) error {
	ctx := c.Request().Context()
	logg := appctx.LoggerFromContext(ctx)

	var req model.TelegramCodeChatID
	if err := c.Bind(&req); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(req); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	err := handler.db.BindTelegramChat(ctx, req.ChatID, req.Code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logg.Warn("not found", logger.ErrAttr(err))

			return c.JSON(http.StatusNotFound, map[string]string{"err": "not found"})
		}

		logg.Error("Database error", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, map[string]string{"err": "Server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"telegram bound": "telegram bound",
	})
}

func (handler *Handler) Verify2FA(c echo.Context) error {
	ctx := c.Request().Context()
	logg := appctx.LoggerFromContext(ctx)

	var req model.Verify2FARequest
	if err := c.Bind(&req); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(req); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", logger.ErrAttr(err))

		return c.JSON(http.StatusBadRequest, errResp)
	}

	tokens, err := handler.service.Verify2FA(&req, ctx, logg)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", logger.ErrAttr(err))

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, tokens)
}
