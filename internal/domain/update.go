package domain

type ScheduleUpdate struct {
	Added   []Event
	Deleted []Event
}

func (u ScheduleUpdate) IsEmpty() bool {
	return len(u.Added) == 0 && len(u.Deleted) == 0
}

func NewScheduleUpdate(old, current Schedule) ScheduleUpdate {
	update := ScheduleUpdate{}

	oldMap := make(map[string]Event)
	for _, e := range old.Events {
		oldMap[e.Hash()] = e
	}

	currMap := make(map[string]Event)
	for _, e := range current.Events {
		currMap[e.Hash()] = e
	}

	// Check what was deleted (was in old, but not in new)
	for _, e := range old.Events {

		if _, ok := currMap[e.Hash()]; !ok {
			update.Deleted = append(update.Deleted, e)
		}
	}

	// Check what was added (is in new, but not in old)
	for _, e := range current.Events {
		if _, ok := oldMap[e.Hash()]; !ok {
			update.Added = append(update.Added, e)
		}
	}

	return update
}
