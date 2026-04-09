package domain

import "fmt"

type EventType uint8

const (
	PreparationEvent = iota
	ProtectionBagrutTestEvent
	BagrutTestEvent
)

func (et EventType) String() string {
	switch et {
	case PreparationEvent:
		return "Tikbur"
	case ProtectionBagrutTestEvent:
		return "Magen"
	case BagrutTestEvent:
		return "Bagrut"
	}
	return "Undefined event"
}

type Event struct {
	Date Date

	Type  EventType
	Group string
	Text  string
}

func NewEvent(date Date, gr string, text string, et EventType) (Event, error) {
	if et > 2 {
		return Event{}, fmt.Errorf("not existing event type")
	}
	return Event{
		Date:  date,
		Group: gr,
		Text:  text,
		Type:  et,
	}, nil
}
