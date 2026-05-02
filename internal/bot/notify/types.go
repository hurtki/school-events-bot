package notify

type Summary struct {
	Rescheduled    []RescheduledEntry `json:"rescheduled"`
	Clarifications []EventEntry       `json:"clarifications"`
	Changes        []ChangeEntry      `json:"changes"`
	Added          []EventEntry       `json:"added"`
	Deleted        []EventEntry       `json:"deleted"`
}

func (s *Summary) IsEmpty() bool {
	return len(s.Rescheduled) == 0 && len(s.Clarifications) == 0 &&
		len(s.Changes) == 0 && len(s.Added) == 0 && len(s.Deleted) == 0
}

type EventEntry struct {
	Date      string `json:"date"`
	Group     string `json:"group"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	Note      string `json:"note,omitempty"`
	SourceURL string `json:"source_url,omitempty"`
}

type ChangeEntry struct {
	Date      string `json:"date"`
	Group     string `json:"group"`
	Type      string `json:"type"`
	OldText   string `json:"old_text"`
	NewText   string `json:"new_text"`
	SourceURL string `json:"source_url,omitempty"`
}

type RescheduledEntry struct {
	OldDate   string `json:"old_date"`
	NewDate   string `json:"new_date"`
	Group     string `json:"group"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	SourceURL string `json:"source_url,omitempty"`
}
