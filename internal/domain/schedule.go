package domain

import (
	"slices"
)

type Schedule struct {
	Events []Event
}

func NewSchedule(evs []Event) (Schedule, error) {
	return Schedule{
		Events: evs,
	}, nil
}

func (s Schedule) GetUpcomingEventsSummary(count int) UpcomingEventsSummary {
	events := make(map[Group][]Event)

	for _, e := range s.Events {
		if e.Date.DaysUntil() < 0 {
			continue
		}
		events[e.Group] = append(events[e.Group], e)
	}

	for gr := range events {
		slices.SortFunc(events[gr], func(a, b Event) int {
			return a.Date.Compare(b.Date)
		})
		if len(events[gr]) > count {
			events[gr] = events[gr][:count]
		}
	}

	return NewUpcomingEventsSummary(events)
}
