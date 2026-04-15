package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/hurtki/school-events-bot/internal/domain"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) SendEventsSummaryAndPin(summary domain.UpcomingEventsSummary) (msgID int, err error) {

	msg, err := b.bot.Send(&tele.Chat{ID: b.cfg.UpdatesChannel},
		b.formatSummary(summary),
		&tele.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeHTML,
		},
	)

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

	_, err := b.bot.Edit(&tele.Message{ID: msgID, Chat: &tele.Chat{ID: b.cfg.UpdatesChannel}},
		b.formatSummary(summary),
		&tele.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeHTML,
		},
	)

	if err != nil {
		return fmt.Errorf("can't edit message: %w", err)
	}
	return nil
}

func (b *Bot) DeleteMessage(msgID int) error {
	err := b.bot.Delete(&tele.Message{ID: msgID, Chat: &tele.Chat{ID: b.cfg.UpdatesChannel}})
	if err != nil {
		return fmt.Errorf("can't delete message: %w", err)
	}
	return nil
}

func (b *Bot) formatSummary(summary domain.UpcomingEventsSummary) string {
	var sb strings.Builder

	groups := []struct {
		label string
		group domain.Group
	}{
		{"📗 שכבת י'", domain.TenthGradeGroup},
		{"📘 שכבת י\"א", domain.EleventhGradeGroup},
		{"📙 שכבת י\"ב", domain.TwelfthGradeGroup},
		{"🎓 מכללה", domain.CollegeGroup},
	}

	sb.WriteString("<b>📅 אירועים חשובים קרובים</b>\n\n")

	for _, g := range groups {
		events := summary.Events[g.group]
		fmt.Fprintf(&sb, "<b>%s</b>\n", g.label)
		if len(events) == 0 {
			sb.WriteString("— אין אירועים\n")
		} else {
			for _, e := range events {
				daysUntilEvent := e.Date.DaysUntil()
				date := ""
				switch daysUntilEvent {
				case 0:
					date = "היום"
				case 1:
					date = "מחר"
				default:
					date = fmt.Sprintf("%d ימים [%s]", daysUntilEvent, e.Date.String())
				}

				if e.SourceURL != "" {
					date = fmt.Sprintf("<a href=\"%s\">%s</a>", e.SourceURL, date)
				}
				fmt.Fprintf(&sb, "\u200F%s\n<b>%s</b>", date, e.Text)
			}
		}
		sb.WriteString("\n")
	}

	loc, _ := time.LoadLocation("Asia/Jerusalem")
	now := time.Now().In(loc)

	fmt.Fprintf(&sb, "<i>עודכן: %02d:%02d:%02d</i>", now.Hour(), now.Minute(), now.Second())

	return sb.String()
}
