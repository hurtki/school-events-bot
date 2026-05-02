package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/config"
	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/hurtki/school-events-bot/internal/infrastructure/ai"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()

	logger.Info("loading config")

	botCfg, err := config.LoadBotConfig(config.EnvFileSource)
	if err != nil {
		logger.Error("can't load bot config", "err", err)
		return
	}
	appCfg, err := config.LoadAppConfig(config.EnvFileSource)
	if err != nil {
		logger.Error("can't load app config", "err", err)
		return
	}
	logger.Info("config loaded")

	var summaryAI interface {
		Text(ctx context.Context, prompt string) (string, error)
	}
	geminiAI, err := ai.NewGeminiAI(appCfg.GeminiAPIKey, appCfg.GeminiModel)
	if err != nil {
		logger.Warn("Gemini unavailable, falling back to basic format", "err", err)
		summaryAI = ai.NewNoopGeminiAI()
	} else {
		logger.Info("Gemini AI initialised", "model", appCfg.GeminiModel)
		summaryAI = geminiAI
	}

	logger.Info("creating bot")
	b, err := bot.NewBot(botCfg, summaryAI, logger)
	if err != nil {
		logger.Error("can't create bot", "err", err)
		return
	}
	logger.Info("bot created")

	mustDate := func(s string) domain.Date {
		d, err := domain.NewDate(s)
		if err != nil {
			panic(err)
		}
		return d
	}

	todayStr := time.Now().Format("2.1.2006")
	tomorrowStr := time.Now().AddDate(0, 0, 1).Format("2.1.2006")

	update := domain.ScheduleUpdate{
		Added: []domain.Event{
			// Clarification: time slot added to existing event
			{
				Date: mustDate("23.4.2026"), Type: domain.ProtectionBagrutTestEvent,
				Group: domain.TwelfthGradeGroup, Text: "מגן לשון עולים ה,ו 08:00-13:00",
				SourceURL: "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/edit#gid=1710319946&range=M319",
			},
			// Clarification: multi-line, teacher+time added at top
			{
				Date: mustDate("20.4.2026"), Type: domain.PreparationEvent,
				Group:     domain.CollegeGroup,
				Text:      "תגבור ניסים 8:00-10:30\nתכן מכני יד1 יד2 1.08-\nתגבור אנה 8:00-10:30 יג 1יג2 1.09-\nתגבור דביר 10:30-13:00 יד1 יד2 1.01-",
				SourceURL: "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/edit#gid=898691425&range=J307",
			},
			// Pure addition
			{
				Date: mustDate("10.5.2026"), Type: domain.BagrutTestEvent,
				Group:     domain.TenthGradeGroup,
				Text:      "בגרות מתמטיקה 5 יח' — כיתות י1, י2",
				SourceURL: "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/edit#gid=111111&range=B5",
			},
			// Rescheduled: same event moved from 1.5.2026 to 8.5.2026
			{
				Date: mustDate("8.5.2026"), Type: domain.BagrutTestEvent,
				Group:     domain.EleventhGradeGroup,
				Text:      "בגרות אנגלית — כיתות יא1, יא2, יא3",
				SourceURL: "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/edit#gid=123456&range=A1",
			},
			// Today's event
			{
				Date:      mustDate(todayStr),
				Type:      domain.ProtectionBagrutTestEvent,
				Group:     domain.TwelfthGradeGroup,
				Text:      "מגן היסטוריה 09:00-12:00",
				SourceURL: "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/edit#gid=222222&range=C8",
			},
			// Tomorrow's event
			{
				Date:      mustDate(tomorrowStr),
				Type:      domain.BagrutTestEvent,
				Group:     domain.TenthGradeGroup,
				Text:      "בגרות ספרות 08:00-11:00",
				SourceURL: "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/edit#gid=333333&range=D12",
			},
		},
		Deleted: []domain.Event{
			// Paired with first added (clarification)
			{
				Date: mustDate("23.4.2026"), Type: domain.ProtectionBagrutTestEvent,
				Group: domain.TwelfthGradeGroup, Text: "מגן לשון עולים ה,ו",
			},
			// Paired with second added (clarification)
			{
				Date: mustDate("20.4.2026"), Type: domain.PreparationEvent,
				Group: domain.CollegeGroup,
				Text:  "כן מכני יד1 יד2 1.08-\nתגבור אנה 8:00-10:30 יג 1יג2 1.09-\nתגבור דביר 10:30-13:00 יד1 יד2 1.01-",
			},
			// Pure deletion
			{
				Date: mustDate("15.5.2026"), Type: domain.ProtectionBagrutTestEvent,
				Group: domain.EleventhGradeGroup, Text: "מגן ביולוגיה",
			},
			// Rescheduled: old date of the moved event
			{
				Date: mustDate("1.5.2026"), Type: domain.BagrutTestEvent,
				Group: domain.EleventhGradeGroup, Text: "בגרות אנגלית — כיתות יא1, יא2, יא3",
			},
		},
	}

	logger.Info("sending test notification",
		"added", len(update.Added),
		"deleted", len(update.Deleted),
	)

	if err := b.NotifyAboutUpdate(ctx, update); err != nil {
		logger.Error("NotifyAboutUpdate failed", "err", err)
		return
	}
	logger.Info("notification sent successfully")
}
