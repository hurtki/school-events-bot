package domain

type UpcomingEventsSummary struct {
	Events map[Group][]Event
}

func NewUpcomingEventsSummary(events map[Group][]Event) UpcomingEventsSummary {
	return UpcomingEventsSummary{
		Events: events,
	}
}
