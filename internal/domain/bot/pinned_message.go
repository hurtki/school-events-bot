package bot

import "time"

type PinnedMessage struct {
	LastMessageID     int
	LastMessageSentAt time.Time
}

func NewPinnedMessage(lastMessageID int, LastMessageSentAt time.Time) (PinnedMessage, error) {
	return PinnedMessage{
		LastMessageID:     lastMessageID,
		LastMessageSentAt: LastMessageSentAt,
	}, nil
}
