package botrunner

import (
	"context"
	"github.com/Kofandr/To-do_list/internal/bothandler"
	"github.com/Kofandr/To-do_list/internal/logger"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

func RunPolling(ctx context.Context, logg *slog.Logger, bot *tgbotapi.BotAPI, apiURL string) error {
	logg.Info("starting telegram bot polling...")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-ctx.Done():
			logg.Info("polling stopped")
			return ctx.Err()
		case update := <-updates:
			if update.Message == nil || !update.Message.IsCommand() {
				continue
			}
			switch update.Message.Command() {
			case "start":
				_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Welcome. To link your account, send /link <code>"))
				if err != nil {
					logg.Error("failed to send /start reply", logger.ErrAttr(err))
				}
			case "link":
				bothandler.HandleLinkCommand(logg, bot, update.Message, apiURL)
			default:
				if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")); err != nil {
					logg.Warn("failed to send Unknown command reply", logger.ErrAttr(err))
				}
			}
		}
	}
}
