package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hurtki/school-events-bot/internal/app/bot/pinned"
	"github.com/hurtki/school-events-bot/internal/app/bot/updates"
	"github.com/hurtki/school-events-bot/internal/app/schedule"
	"github.com/hurtki/school-events-bot/internal/app/workers"
	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/config"
	"github.com/hurtki/school-events-bot/internal/evbus"
	"github.com/hurtki/school-events-bot/internal/infrastructure/spreadsheets"
	repository "github.com/hurtki/school-events-bot/internal/repository/schedule"
)

func main() {
	envSource := config.EnvFileSource

	appCfg, err := config.LoadAppConfig(envSource)
	if err != nil {
		fmt.Println("can't load app config:", err)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Info("starting service")

	botCfg, err := config.LoadBotConfig(envSource)
	if err != nil {
		logger.Error("can't load bot config", "err", err)
		return
	}

	docFetcher := spreadsheets.NewDocsFetcher(http.DefaultClient)

	bot, err := bot.NewBot(botCfg)
	if err != nil {
		logger.Error("can't init bot", "err", err)
		return
	}

	scheduleUpdateEventBus := evbus.NewScheduleUpdateEventBus()

	scheduleRepo := repository.NewFileScheduleRepository(appCfg.JsonScheduleFileRepositoryPath)

	scheduleService := schedule.NewScheduleService(docFetcher, appCfg.SpreadsheetsDocumentID, scheduleRepo, scheduleUpdateEventBus)

	botScheduleUpdatesService := updates.NewBotScheduleUpdatesService(logger, bot)

	botUpcomingEventsPinService := pinned.NewBotUpcomingEventsPinService(
		logger,
		nil,
		scheduleRepo,
	)

	scheduleUpdateEventBus.Subscribe(botScheduleUpdatesService.HandleScheduleUpdate)
	scheduleUpdateEventBus.Subscribe(botUpcomingEventsPinService.HandleScheduleUpdate)

	scheduleWorker := workers.NewScheduleWorker(
		logger,
		scheduleService,
		appCfg.SchedulePollerInterval,
	)
	upcomingEventsWorker := workers.NewUpcomingEventsWorker(logger, botUpcomingEventsPinService, appCfg.SchedulePollerInterval)

	scheduleWorker.Start()
	upcomingEventsWorker.Start()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	quitCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	scheduleWorker.Close(quitCtx)
	upcomingEventsWorker.Close(quitCtx)

	cancel()

	// main1(docFetcher, appCfg)
}
