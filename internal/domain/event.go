package domain

import "fmt"

type EventType uint8

const (
	PreparationEvent = iota
	ProtectionBagrutTestEvent
	BagrutTestEvent
	ExamEvent
)

func (et EventType) String() string {
	switch et {
	case PreparationEvent:
		return "Tikbur"
	case ProtectionBagrutTestEvent:
		return "Magen"
	case BagrutTestEvent:
		return "Bagrut"
	case ExamEvent:
		return "Mivhan"
	}
	return "Undefined event"
}

type Event struct {
	Date Date

	Type  EventType
	Group Group
	Text  string

	SourceURL string
}

func NewEvent(date Date, gr Group, text string, et EventType, sourceURL string) (Event, error) {
	if et > 3 {
		return Event{}, fmt.Errorf("not existing event type")
	}
	return Event{
		Date:      date,
		Group:     gr,
		Text:      text,
		Type:      et,
		SourceURL: sourceURL,
	}, nil
}

// Hash is used to compare two events, if they are identical
func (e Event) Hash() string {
	return fmt.Sprintf("%s%s%s%s",
		e.Date.String(),
		e.Type.String(),
		e.Group.String(),
		e.Text)
}
