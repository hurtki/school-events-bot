package bot

import (
	"context"
	"log/slog"
	"time"

	"github.com/hurtki/school-events-bot/internal/config"
	tele "gopkg.in/telebot.v4"
)

type summaryAI interface {
	Text(ctx context.Context, prompt string) (string, error)
}

type Bot struct {
	bot    *tele.Bot
	logger *slog.Logger
	cfg    config.BotConfig
	ai     summaryAI
}

func NewBot(cfg config.BotConfig, ai summaryAI, logger *slog.Logger) (*Bot, error) {
	set := tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: time.Second * 10},
	}
	b, err := tele.NewBot(set)
	if err != nil {
		return nil, err
	}
	return &Bot{
		bot:    b,
		cfg:    cfg,
		logger: logger.With("service", "bot-infra"),
		ai:     ai,
	}, nil
}

func (b *Bot) Close() {
	b.bot.Stop()
}

func (b *Bot) Start() {
	go b.bot.Start()
}
