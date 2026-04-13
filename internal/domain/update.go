package domain

func CompareSchedules(old, current Schedule) ScheduleUpdate {
	update := ScheduleUpdate{}

	oldMap := make(map[Event]bool)
	for _, e := range old.Events {
		oldMap[e] = true
	}

	currMap := make(map[Event]bool)
	for _, e := range current.Events {
		currMap[e] = true
	}

	// Check what was deleted (was in old, but not in new)
	for _, e := range old.Events {
		if !currMap[e] {
			update.Deleted = append(update.Deleted, e)
		}
	}

	// Check what was added (is in new, but not in old)
	for _, e := range current.Events {
		if !oldMap[e] {
			update.Added = append(update.Added, e)
		}
	}

	return update
}

type ScheduleUpdate struct {
	Added   []Event
	Deleted []Event
}

func (u ScheduleUpdate) IsEmpty() bool {
	return len(u.Added) == 0 && len(u.Deleted) == 0
}

type ScheduleUpdateEvent struct {
	Update ScheduleUpdate
}
