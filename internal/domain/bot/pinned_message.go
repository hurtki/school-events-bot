package bot

import "time"

type PinnedMessage struct {
	LastMessageID     int
	LastMessageSentAt time.Time
}
