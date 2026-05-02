package bot

import (
	"context"
	"time"

	"github.com/hurtki/school-events-bot/internal/bot/notify"
	"github.com/hurtki/school-events-bot/internal/domain"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) NotifyAboutUpdate(ctx context.Context, update domain.ScheduleUpdate) error {
	if update.IsEmpty() {
		return nil
	}

	today := time.Now().Format("2.1.2006")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2.1.2006")

	summary, err := notify.GenerateSummary(ctx, b.ai, update, today, tomorrow)
	var msg string
	if err != nil {
		b.logger.Warn("AI summary failed, using basic format", "err", err)
		msg = notify.FormatFallback(update)
	} else if summary.IsEmpty() {
		b.logger.Info("AI determined no meaningful changes, skipping notification")
		return nil
	} else {
		msg = notify.Format(summary, today, tomorrow)
	}

	_, err = b.bot.Send(&tele.Chat{ID: b.cfg.UpdatesChannel}, msg, &tele.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tele.ModeHTML,
	})
	return err
}
