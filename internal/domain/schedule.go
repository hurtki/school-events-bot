package domain

type Schedule struct {
	Events []Event
}

func NewSchedule(evs []Event) (Schedule, error) {
	return Schedule{
		Events: evs,
	}, nil
}
