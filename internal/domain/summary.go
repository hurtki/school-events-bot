package domain

type UpcomingEventsSummary struct {
	events map[Group][]Event
}

func NewUpcomingEventsSummary(events map[Group][]Event) UpcomingEventsSummary {
	return UpcomingEventsSummary{
		events: events,
	}
}
