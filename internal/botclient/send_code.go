package botclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/logger"
	"log/slog"
	"net/http"
)

func SendCode(botURL string, chatID int64, code string, logg *slog.Logger) {
	sendURL := fmt.Sprintf("%s/send-code", botURL)

	reqBody := model.TelegramCodeChatID{
		ChatID: chatID,
		Code:   code,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		logg.Error("failed to marshal send-code request", logger.ErrAttr(err))
		return
	}

	resp, err := http.Post(sendURL, "application/json", bytes.NewReader(body))
	if err != nil {
		logg.Error("failed to call bot /send-code", logger.ErrAttr(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logg.Error("bot /send-code returned non-200", "status", resp.StatusCode)
	}
}
