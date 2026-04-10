package bot

import (
	"context"
	"time"

	"github.com/hurtki/school-events-bot/internal/config"
	tele "github.com/tucnak/telebot"
)

type Bot struct {
	bot *tele.Bot

	cfg config.BotConfig
}

func NewBot(cfg config.BotConfig) (*Bot, error) {
	set := tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: time.Second * 10},
	}
	b, err := tele.NewBot(set)
	if err != nil {
		return nil, err
	}
	return &Bot{
		bot: b,
		cfg: cfg,
	}, nil
}

func (b *Bot) NotifyAboutUpdate(ctx context.Context) error {
	msg := `
	<a href="https://t.me/PoiskEng_bot">Отправить пост</a>
	`
	_, err := b.bot.Send(&tele.Chat{ID: b.cfg.UpdatesChannel}, msg, tele.ModeHTML)
	return err
}
