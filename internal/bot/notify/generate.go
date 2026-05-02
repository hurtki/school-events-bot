package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hurtki/school-events-bot/internal/domain"
)

type AIClient interface {
	Text(ctx context.Context, prompt string) (string, error)
}

type inputEvent struct {
	Date      string `json:"date"`
	Type      string `json:"type"`
	Group     string `json:"group"`
	Text      string `json:"text"`
	SourceURL string `json:"source_url,omitempty"`
}

type inputPayload struct {
	Added   []inputEvent `json:"added"`
	Deleted []inputEvent `json:"deleted"`
}

func GenerateSummary(ctx context.Context, ai AIClient, update domain.ScheduleUpdate, today, tomorrow string) (*Summary, error) {
	p := inputPayload{}
	for _, e := range update.Added {
		p.Added = append(p.Added, inputEvent{
			Date:      e.Date.String(),
			Type:      e.Type.String(),
			Group:     e.Group.String(),
			Text:      e.Text,
			SourceURL: e.SourceURL,
		})
	}
	for _, e := range update.Deleted {
		p.Deleted = append(p.Deleted, inputEvent{
			Date:  e.Date.String(),
			Type:  e.Type.String(),
			Group: e.Group.String(),
			Text:  e.Text,
		})
	}

	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshal input: %w", err)
	}

	raw, err := ai.Text(ctx, buildPrompt(string(data), today, tomorrow))
	if err != nil {
		return nil, err
	}

	cleaned := cleanJSON(raw)
	var summary Summary
	if err := json.Unmarshal([]byte(cleaned), &summary); err != nil {
		return nil, fmt.Errorf("unmarshal AI response: %w", err)
	}
	fillSourceURLs(&summary, update)
	return &summary, nil
}

// fillSourceURLs restores source_url for every summary entry from the original
// update, so we don't rely on the AI copying it correctly.
func fillSourceURLs(s *Summary, update domain.ScheduleUpdate) {
	// key: "date|type|group|text" → sourceURL (only added events have URLs)
	m := make(map[string]string, len(update.Added))
	for _, e := range update.Added {
		if e.SourceURL != "" {
			m[e.Date.String()+"|"+e.Type.String()+"|"+e.Group.String()+"|"+e.Text] = e.SourceURL
		}
	}
	lookup := func(date, typ, group, text string) string {
		return m[date+"|"+typ+"|"+group+"|"+text]
	}

	for i := range s.Clarifications {
		s.Clarifications[i].SourceURL = lookup(s.Clarifications[i].Date, s.Clarifications[i].Type, s.Clarifications[i].Group, s.Clarifications[i].Text)
	}
	for i := range s.Changes {
		s.Changes[i].SourceURL = lookup(s.Changes[i].Date, s.Changes[i].Type, s.Changes[i].Group, s.Changes[i].NewText)
	}
	for i := range s.Added {
		s.Added[i].SourceURL = lookup(s.Added[i].Date, s.Added[i].Type, s.Added[i].Group, s.Added[i].Text)
	}
	for i := range s.Rescheduled {
		s.Rescheduled[i].SourceURL = lookup(s.Rescheduled[i].NewDate, s.Rescheduled[i].Type, s.Rescheduled[i].Group, s.Rescheduled[i].Text)
	}
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	for _, fence := range []string{"```json", "```"} {
		if strings.HasPrefix(s, fence) {
			s = strings.TrimPrefix(s, fence)
			s = strings.TrimSuffix(s, "```")
			s = strings.TrimSpace(s)
			break
		}
	}
	return s
}
