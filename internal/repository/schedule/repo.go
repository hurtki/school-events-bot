package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/hurtki/school-events-bot/internal/domain"
)

type FileScheduleRepository struct {
	path string
	mu   sync.RWMutex
}

func NewFileScheduleRepository(path string) *FileScheduleRepository {
	return &FileScheduleRepository{
		path: path,
	}
}

func (r *FileScheduleRepository) GetLastSchedule(ctx context.Context) (*domain.Schedule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.Open(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var schedule domain.Schedule
	if err := json.NewDecoder(file).Decode(&schedule); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	return &schedule, nil
}

func (r *FileScheduleRepository) SaveSchedule(ctx context.Context, s domain.Schedule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Create(r.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}
