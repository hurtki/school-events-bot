package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/hurtki/school-events-bot/internal/config"
	"github.com/hurtki/school-events-bot/internal/infrastructure/spreadsheets"
)

func main() {
	test()
	return
	envSource := config.EnvFileSource

	appCfg, err := config.LoadAppConfig(envSource)
	if err != nil {
		fmt.Println("can't load config", err)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Info("starting service")

	// botCfg, err := config.LoadBotConfig(envSource)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	docFetcher := spreadsheets.NewDocsFetcher(http.DefaultClient)

	// bot, err := bot.NewBot(botCfg)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// scheduleService := schedule.NewScheduleService(docFetcher, appCfg.SpreadsheetsDocumentID)
	//
	// scheduleRepo := repository.NewFileScheduleRepository(appCfg.JsonScheduleFileRepositoryPath)

	// poller := workers.NewSchedulePoller(
	// 	logger,
	// 	scheduleService,
	// 	bot,
	// 	appCfg.SchedulePollerInterval,
	// 	scheduleRepo,
	// )
	// poller.Start()
	// // graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit
	// quitCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// poller.Close(quitCtx)

	main1(docFetcher, appCfg)
	// cancel()
}
