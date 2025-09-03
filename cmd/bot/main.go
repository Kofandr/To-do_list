package main

import (
	"context"
	"errors"
	"github.com/Kofandr/To-do_list/config/botconfig"
	"github.com/Kofandr/To-do_list/internal/bothandler"
	"github.com/Kofandr/To-do_list/internal/botrunner"
	"github.com/Kofandr/To-do_list/internal/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	var exitCode int
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	cfg := botconfig.MustLoad()
	logg := logger.New(cfg.LoggerLevel)
	slog.SetDefault(logg)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logg.Error("failed to create Telegram bot", logger.ErrAttr(err))
		os.Exit(1)
	}
	logg.Info("bot authorized", "username", bot.Self.UserName)

	e := echo.New()
	e.POST("/send-code", bothandler.SendCodeHandler(logg, bot))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	errCh := make(chan error, 2)

	pollCtx, pollCancel := context.WithCancel(context.Background())
	defer pollCancel()
	go func() {
		if err := botrunner.RunPolling(pollCtx, logg, bot, cfg.APIURL); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
	}()

	go func() {
		logg.Info("starting bot HTTP server", "port", cfg.BotPort)
		if err := e.Start(":" + cfg.BotPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	var runErr error
	select {
	case sig := <-sigCh:
		logg.Info("shutdown requested", "signal", sig.String())
	case runErr = <-errCh:
		logg.Error("bot error", logger.ErrAttr(runErr))
	}

	pollCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logg.Error("graceful shutdown failed", logger.ErrAttr(err))
	} else {
		logg.Info("graceful shutdown completed")
	}

	if runErr != nil {
		exitCode = 1
	}
}
