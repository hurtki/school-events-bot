package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/hurtki/school-events-bot/internal/domain"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) NotifyAboutUpdate(ctx context.Context, update domain.ScheduleUpdate) error {
	if update.IsEmpty() {
		return nil
	}

	var sb strings.Builder

	// Header
	sb.WriteString("<b>🔔 עדכון במערכת!</b>\n\n")

	// Helper to get Hebrew group name with emoji
	getGroupLabel := func(g domain.Group) string {
		switch g {
		case domain.TenthGradeGroup:
			return "📗 שכבת י'"
		case domain.EleventhGradeGroup:
			return "📘 שכבת י\"א"
		case domain.TwelfthGradeGroup:
			return "📙 שכבת י\"ב"
		case domain.CollegeGroup:
			return "🎓 מכללה"
		default:
			return "👥"
		}
	}

	// Added Events
	if len(update.Added) > 0 {
		sb.WriteString("<b>🆕 אירועים חדשים שנוספו:</b>\n")
		for _, e := range update.Added {
			groupLabel := getGroupLabel(e.Group)

			date := fmt.Sprintf("[%s]", e.Date.String())
			if e.SourceURL != "" {
				date = fmt.Sprintf("<a href=\"%s\">%s</a>", e.SourceURL, date)
			}

			text := fmt.Sprintf("<b>%s</b>", e.Text)

			fmt.Fprintf(&sb, "• %s | %s\n%s", groupLabel, date, text)
		}
		sb.WriteString("\n")
	}

	// Deleted Events
	if len(update.Deleted) > 0 {
		sb.WriteString("<b>❌ אירועים שנמחקו:</b>\n")
		for _, e := range update.Deleted {
			groupLabel := getGroupLabel(e.Group)

			date := fmt.Sprintf("[%s]", e.Date.String())
			text := fmt.Sprintf("<s>%s</s>", e.Text)

			fmt.Fprintf(&sb, "• %s | %s\n%s", groupLabel, date, text)
		}
	}

	_, err := b.bot.Send(&tele.Chat{ID: b.cfg.UpdatesChannel}, sb.String(), &tele.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tele.ModeHTML,
	})
	return err
}
