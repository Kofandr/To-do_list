package bothandler

import (
	"bytes"
	"encoding/json"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/logger"
	"log/slog"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleLinkCommand(logg *slog.Logger, bot *tgbotapi.BotAPI, message *tgbotapi.Message, apiURL string) {
	code := strings.TrimSpace(message.CommandArguments())
	if code == "" {
		_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Please provide a link code after /link"))
		if err != nil {
			logg.Error("failed to send /link message", logger.ErrAttr(err))
		}
		return
	}

	request := model.TelegramCodeChatID{
		ChatID: message.Chat.ID,
		Code:   code,
	}

	body, err := json.Marshal(request)
	if err != nil {
		logg.Error("failed to marshal link request", logger.ErrAttr(err))
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Internal error"))
		return
	}

	resp, err := http.Post(apiURL+"/2fa/bind", "application/json", bytes.NewBuffer(body))
	if err != nil {
		logg.Error("failed to call /2fa/bind", logger.ErrAttr(err))
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Failed to reach server"))
		return
	}
	defer resp.Body.Close()

	masked := maskCode(code)
	logg = logg.With("code", masked, "chat_id", message.Chat.ID)

	if resp.StatusCode != http.StatusOK {
		logg.Warn("invalid link code")

		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Invalid or expired code"))
		return
	}

	logg.Info("telegram linked")
	_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Telegram linked successfully"))
}

func maskCode(code string) string {
	if len(code) <= 2 {
		return "**"
	}
	return strings.Repeat("*", len(code)-2) + code[len(code)-2:]
}
