package bothandler

import (
	"fmt"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func SendCodeHandler(logg *slog.Logger, bot *tgbotapi.BotAPI) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.TelegramCodeChatID
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		}

		msg := tgbotapi.NewMessage(req.ChatID, fmt.Sprintf("Your 2FA code: %s", req.Code))
		if _, err := bot.Send(msg); err != nil {
			logg.Error("failed to send telegram message", logger.ErrAttr(err))
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "send failed"})
		}

		return c.JSON(http.StatusOK, model.SuccessResponse{Message: "sent"})
	}
}
