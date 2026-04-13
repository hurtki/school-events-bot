package domain

import (
	"slices"
	"time"
)

type Schedule struct {
	Events []Event
}

func NewSchedule(evs []Event) (Schedule, error) {
	return Schedule{
		Events: evs,
	}, nil
}

// GetUpcomingEventsSummary returns the nearest events for each group,
// prioritizing the specified priorityTypes until the "first" limit is reached.
// If priority events are insufficient, it fills the remaining slots with other upcoming events.
func (s Schedule) GetUpcomingEventsSummary(first int, priorityTypes ...EventType) UpcomingEventsSummary {
	prs := make(map[EventType]struct{})
	for _, p := range priorityTypes {
		prs[p] = struct{}{}
	}

	events := make(map[Group][]Event)

	for _, e := range s.Events {
		if e.Date.T.Before(time.Now()) {
			continue
		}
		events[e.Group] = append(events[e.Group], e)
	}

	for gr := range events {
		slices.SortFunc(events[gr], func(a Event, b Event) int {
			return a.Date.Compare(b.Date)
		})

		grEvs := make([]Event, 0, first)

		for _, e := range events[gr] {
			if len(grEvs) >= first {
				break
			}
			if _, ok := prs[e.Type]; ok {
				grEvs = append(grEvs, e)
			}
		}
		if len(grEvs) < first {
			for _, e := range events[gr] {
				if len(grEvs) >= first {
					break
				}
				if _, ok := prs[e.Type]; !ok {
					grEvs = append(grEvs, e)
				}
			}
		}
		events[gr] = grEvs
		slices.SortFunc(events[gr], func(a Event, b Event) int {
			return a.Date.Compare(b.Date)
		})
	}

	return NewUpcomingEventsSummary(events)
}
