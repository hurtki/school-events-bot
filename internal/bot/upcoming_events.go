package bot

import (
	"fmt"
	"time"

	"github.com/hurtki/school-events-bot/internal/domain"
	tele "github.com/tucnak/telebot"
)

func (b *Bot) SendEventsSummaryAndPin(summary domain.UpcomingEventsSummary) (msgID int, err error) {

	msg, err := b.bot.Send(&tele.Chat{ID: b.cfg.UpdatesChannel}, b.formatSummary(summary))

	if err != nil {
		return 0, fmt.Errorf("can't send message: %w", err)
	}

	err = b.bot.Pin(msg, tele.Silent)
	if err != nil {
		b.logger.Error("can't pin message", "err", err)
		return msgID, nil
	}

	return msg.ID, nil
}

func (b *Bot) UpdateMessageWithEventsSummary(msgID int, summary domain.UpcomingEventsSummary) error {

	_, err := b.bot.Edit(&tele.Message{ID: msgID, Chat: &tele.Chat{ID: b.cfg.UpdatesChannel}}, b.formatSummary(summary))
	if err != nil {
		return fmt.Errorf("can't edit message: %w", err)
	}
	return nil
}

func (b *Bot) DeleteMessageWithSummary(msgID int) error {
	err := b.bot.Delete(&tele.Message{ID: msgID, Chat: &tele.Chat{ID: b.cfg.UpdatesChannel}})
	if err != nil {
		return fmt.Errorf("can't delete message: %w", err)
	}
	return nil
}

func (b *Bot) formatSummary(summary domain.UpcomingEventsSummary) string {
	return fmt.Sprintf(`
	10:
	%+v
	11:
	%+v
	12:
	%+v
	college:
	%+v

	updated at: %d:%d:%d
	`, summary.Events[domain.TenthGradeGroup], summary.Events[domain.EleventhGradeGroup], summary.Events[domain.TwelfthGradeGroup], summary.Events[domain.CollegeGroup], time.Now().Hour(), time.Now().Minute(), time.Now().Second())
}
