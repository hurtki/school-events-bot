package repository

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	bot_domain "github.com/hurtki/school-events-bot/internal/domain/bot"
)

type JSONPinnedMessageRepo struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONPinnedMessageRepo(filePath string) *JSONPinnedMessageRepo {
	return &JSONPinnedMessageRepo{
		filePath: filePath,
	}
}

func (r *JSONPinnedMessageRepo) GetLastPinnedMessage(ctx context.Context) (bot_domain.PinnedMessage, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return bot_domain.PinnedMessage{}, false, nil
	}

	var msg bot_domain.PinnedMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return bot_domain.PinnedMessage{}, false, nil
	}

	return msg, true, nil
}

func (r *JSONPinnedMessageRepo) SaveLastPinnedMessage(ctx context.Context, msg bot_domain.PinnedMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	tmpFile := r.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, r.filePath)
}
