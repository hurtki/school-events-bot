package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hurtki/school-events-bot/internal/app/schedule"
	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/config"
	"github.com/hurtki/school-events-bot/internal/infrastructure/spreadsheets"
)

func main() {
	envSource := config.EnvFileSource

	appCfg, err := config.LoadAppConfig(envSource)
	if err != nil {
		fmt.Println(err)
		return
	}

	botCfg, err := config.LoadBotConfig(envSource)
	if err != nil {
		fmt.Println(err)
		return
	}

	docFetcher := spreadsheets.NewDocsFetcher(appCfg.SpreadsheetsDocumentID, http.DefaultClient)

	bot, err := bot.NewBot(botCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = bot.NotifyAboutUpdate(context.Background())

	scheduleService := schedule.NewScheduleService(docFetcher)

	_, _ = scheduleService.GetSchedule(context.Background())

	main1(docFetcher)
}
