package notify

import (
	"fmt"
	"strings"

	"github.com/hurtki/school-events-bot/internal/domain"
)

const rtlMark = "‏"

func Format(s *Summary, today, tomorrow string) string {
	var sb strings.Builder
	sb.WriteString("<b>🔔 עדכון במערכת!</b>")

	if len(s.Clarifications) > 0 {
		sb.WriteString("\n\n<b>📌 הבהרות:</b>\n")
		for _, e := range s.Clarifications {
			sb.WriteString("• ")
			sb.WriteString(groupLabel(e.Group))
			sb.WriteString(" | ")
			sb.WriteString(dateLabel(e.Date, today, tomorrow, e.SourceURL))
			if e.Note != "" {
				sb.WriteString(" (")
				sb.WriteString(e.Note)
				sb.WriteString(")")
			}
			sb.WriteString("\n")
			sb.WriteString(e.Text)
		}
	}

	if len(s.Rescheduled) > 0 {
		sb.WriteString("\n\n<b>📅 שינוי תאריך:</b>\n")
		for _, r := range s.Rescheduled {
			sb.WriteString("• ")
			sb.WriteString(groupLabel(r.Group))
			sb.WriteString(" |")
			sb.WriteString(r.OldDate)
			sb.WriteString("←")
			sb.WriteString(dateLabel(r.NewDate, today, tomorrow, r.SourceURL))
			sb.WriteString("\n")
			sb.WriteString(r.Text)
		}
	}

	if len(s.Changes) > 0 {
		sb.WriteString("\n\n<b>✏️ שינויים:</b>\n")
		for _, c := range s.Changes {
			sb.WriteString("• ")
			sb.WriteString(groupLabel(c.Group))
			sb.WriteString(" | ")
			sb.WriteString(dateLabel(c.Date, today, tomorrow, c.SourceURL))
			sb.WriteString("\n")
			sb.WriteString(renderChange(c))
		}
	}

	if len(s.Added) > 0 {
		sb.WriteString("\n\n<b>🆕 אירועים חדשים שנוספו:</b>\n")
		for _, e := range s.Added {
			sb.WriteString("• ")
			sb.WriteString(groupLabel(e.Group))
			sb.WriteString(" | ")
			sb.WriteString(dateLabel(e.Date, today, tomorrow, e.SourceURL))
			sb.WriteString("\n")
			sb.WriteString(e.Text)
		}
	}

	if len(s.Deleted) > 0 {
		sb.WriteString("\n\n<b>❌ אירועים שנמחקו:</b>\n")
		for _, e := range s.Deleted {
			sb.WriteString("• ")
			sb.WriteString(groupLabel(e.Group))
			sb.WriteString(" | ")
			sb.WriteString(dateLabel(e.Date, today, tomorrow, ""))
			sb.WriteString("\n<s>")
			sb.WriteString(e.Text)
			sb.WriteString("</s>")
		}
	}

	return addRTLMarks(truncate(sb.String(), 4096))
}

func FormatFallback(update domain.ScheduleUpdate) string {
	var sb strings.Builder
	sb.WriteString("<b>🔔 עדכון במערכת!</b>")

	if len(update.Added) > 0 {
		sb.WriteString("\n\n<b>🆕 אירועים חדשים שנוספו:</b>")
		for _, e := range update.Added {
			dl := dateLabel(e.Date.String(), "", "", e.SourceURL)
			fmt.Fprintf(&sb, "\n• %s | %s\n%s", groupLabel(e.Group.String()), dl, e.Text)
		}
	}

	if len(update.Deleted) > 0 {
		sb.WriteString("\n\n<b>❌ אירועים שנמחקו:</b>")
		for _, e := range update.Deleted {
			dl := dateLabel(e.Date.String(), "", "", "")
			fmt.Fprintf(&sb, "\n• %s | %s\n<s>%s</s>", groupLabel(e.Group.String()), dl, e.Text)
		}
	}

	return addRTLMarks(sb.String())
}

func renderChange(c ChangeEntry) string {
	segs := computeDiff(c.OldText, c.NewText)
	var sb strings.Builder

	if len(segs) == 0 {
		// No diff detected — show new text only.
		sb.WriteString(c.NewText)
		return sb.String()
	}

	for _, seg := range segs {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		switch {
		case seg.ellipsis:
			sb.WriteString("...")
		case seg.removed:
			sb.WriteString("<s>")
			sb.WriteString(seg.text)
			sb.WriteString("</s>")
		}
	}
	sb.WriteString("\n")
	sb.WriteString(c.NewText)
	return sb.String()
}

func dateLabel(date, today, tomorrow, sourceURL string) string {
	var label string
	switch date {
	case today:
		label = "❗ היום"
	case tomorrow:
		label = "⚠️ מחר"
	default:
		label = date
	}
	if sourceURL != "" {
		return fmt.Sprintf(`<a href="%s">%s</a>`, sourceURL, label)
	}
	return label
}

func groupLabel(g string) string {
	switch g {
	case "10th Grade":
		return "📗 שכבת י'"
	case "11th Grade":
		return `📘 שכבת י"א`
	case "12th Grade":
		return `📙 שכבת י"ב`
	case "College":
		return "🎓 מכללה"
	default:
		return "👥"
	}
}

func addRTLMarks(s string) string {
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if !strings.HasPrefix(line, rtlMark) {
			lines[i] = rtlMark + line
		}
	}
	return strings.Join(lines, "\n")
}

const truncationSuffix = "\n<i>...ההודעה קוצרה</i>"

func truncate(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	suffix := []rune(truncationSuffix)
	return string(runes[:maxRunes-len(suffix)]) + truncationSuffix
}
