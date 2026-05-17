package parser

import (
	"strings"

	"github.com/hurtki/school-events-bot/internal/domain"
)

const (
	preparation          = "תגבור"
	protectionBagrutTest = "מגן"
	bagrutTest           = "בגרות"
	mivhanTest           = "מבחן"
)

var (
	baseSpreadsheetURL = "https://docs.google.com/spreadsheets/d/"
)

// parseDayIntoEvents gets all entries for one day.
// Contains core logic of how to separate and compose events.
func parseDayIntoEvents(
	day string,
	group domain.Group,
	date domain.Date,
	daySrcURL string,
) []domain.Event {

	// we accumulate text and event type of events
	// for example
	// we get new line where written that there is an event
	// we start to add every new line to it, but
	// if we get line with other type, then we create event
	// with accumulated data and start accumulate data for new one
	var events []domain.Event
	var text string
	var et domain.EventType

	// "saves" accumulated data ( if there is )
	// if new type is other than was before
	flushIfType := func(incoming domain.EventType) {
		if text != "" && et != incoming {
			if event, err := domain.NewEvent(date, group, text, et, daySrcURL); err == nil {
				events = append(events, event)
			}
			text = ""
		}
	}

	// "saves" accumulated data ( if there is ) strictly
	flush := func() {
		if text != "" {
			if event, err := domain.NewEvent(date, group, text, et, daySrcURL); err == nil {
				events = append(events, event)
			}
			text = ""
		}
	}

	for line := range strings.SplitSeq(day, "\n") {
		switch {
		case strings.Contains(line, preparation):
			flushIfType(domain.PreparationEvent)
			text += line + "\n"
			et = domain.PreparationEvent

		case strings.Contains(line, protectionBagrutTest):
			flushIfType(domain.ProtectionBagrutTestEvent)
			text += line + "\n"
			et = domain.ProtectionBagrutTestEvent

		case strings.Contains(line, bagrutTest):
			flushIfType(domain.BagrutTestEvent)
			text += line + "\n"
			et = domain.BagrutTestEvent

		case strings.Contains(line, mivhanTest):
			flushIfType(domain.ExamEvent)
			text += line + "\n"
			et = domain.ExamEvent

		default:
			if strings.TrimSpace(line) == "" {
				flush()
			} else if text != "" {
				text += line + "\n"
			}
		}
	}
	flush()

	return events
}
