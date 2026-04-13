package domain

import "time"

type Schedule struct {
	Events []Event
}

func NewSchedule(evs []Event) (Schedule, error) {
	return Schedule{
		Events: evs,
	}, nil
}

// GetUpcomingEventsSummary exports upcoming events for every group
// If neededTypes were given, they will be returned, in other case all types of events will be returned
func (s Schedule) GetUpcomingEventsSummary(first int, neededTypes ...EventType) UpcomingEventsSummary {
	events := make(map[Group][]Event)

	neededTypesFat := func(neededTypes []EventType) map[EventType]struct{} {
		res := make(map[EventType]struct{})
		for _, nt := range neededTypes {
			res[nt] = struct{}{}
		}
		return res
	}(neededTypes)

	for _, e := range s.Events {
		// check if type's event is needed

		// if event's date is before now, it's not upcoming
		if e.Date.T.Before(time.Now()) {
			continue
		}

		evs, ok := events[e.Group]
		if !ok {
			events[e.Group] = []Event{e}
			continue
		}

		if len(evs) < first {
			evs = append(evs, e)
			events[e.Group] = evs
			continue
		}

		// if some events is before or at the same date, ignoring it
		if evs[len(evs)-1].Date.T.Compare(e.Date.T) < 1 {
			continue
		}

		for i := len(evs) - 1; i >= 0; i-- {
			existingE := evs[i]

		}

	}
}
