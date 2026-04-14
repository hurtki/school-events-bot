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
	pinned_message_repo "github.com/hurtki/school-events-bot/internal/repository/pinned_message"
	schedule_repo "github.com/hurtki/school-events-bot/internal/repository/schedule"
)

func main() {
	envSource := config.EnviromentVariablesSource

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

	bot, err := bot.NewBot(botCfg, logger)
	if err != nil {
		logger.Error("can't init bot", "err", err)
		return
	}

	scheduleUpdateEventBus := evbus.NewScheduleUpdateEventBus()

	scheduleRepo := schedule_repo.NewFileScheduleRepository(appCfg.JsonScheduleFileRepositoryPath)
	pinnedMsgStateRepo := pinned_message_repo.NewJSONPinnedMessageRepo(appCfg.JsonPinnedMessageStateFileRepositoryPath)

	scheduleService := schedule.NewScheduleService(docFetcher, appCfg.SpreadsheetsDocumentID, scheduleRepo, scheduleUpdateEventBus)

	botScheduleUpdatesService := updates.NewBotScheduleUpdatesService(logger, bot)

	botUpcomingEventsPinService := pinned.NewBotUpcomingEventsPinService(
		logger,
		pinnedMsgStateRepo,
		scheduleRepo,
		bot,
	)

	scheduleUpdateEventBus.Subscribe(botScheduleUpdatesService.HandleScheduleUpdate)
	scheduleUpdateEventBus.Subscribe(botUpcomingEventsPinService.HandleScheduleUpdate)

	scheduleWorker := workers.NewScheduleWorker(
		logger,
		scheduleService,
		appCfg.ScheduleWorkerInterval,
	)
	upcomingEventsWorker := workers.NewUpcomingEventsWorker(logger, botUpcomingEventsPinService, appCfg.UpcomingEventsWorkerInterval)

	bot.DeletePinsHandle()

	scheduleWorker.Start()
	upcomingEventsWorker.Start()
	bot.Start()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	quitCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	scheduleWorker.Close(quitCtx)
	upcomingEventsWorker.Close(quitCtx)
	bot.Close()

	cancel()

	// main1(docFetcher, appCfg)
}
