package schedule

import (
	"context"
	"fmt"
	"io"

	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/hurtki/school-events-bot/internal/parser"
)

type XLSXScheduleDocumentFetcher interface {
	FetchXLSX(ctx context.Context, docID string) (io.ReadCloser, error)
}

type ScheduleRepository interface {
	GetLastSchedule(ctx context.Context) (*domain.Schedule, error)
	SaveSchedule(ctx context.Context, s domain.Schedule) error
}

type ScheduleUpdatesEventBus interface {
	Publish(context.Context, domain.ScheduleUpdate)
}

type ScheduleService struct {
	fetcher XLSXScheduleDocumentFetcher
	repo    ScheduleRepository
	docID   string
	evBus   ScheduleUpdatesEventBus
}

func NewScheduleService(fetcher XLSXScheduleDocumentFetcher, docID string, repo ScheduleRepository, evBus ScheduleUpdatesEventBus) *ScheduleService {
	return &ScheduleService{
		fetcher: fetcher,
		docID:   docID,
		repo:    repo,
		evBus:   evBus,
	}
}

func (s *ScheduleService) Update(ctx context.Context) error {
	xlsx, err := s.fetcher.FetchXLSX(ctx, s.docID)
	if err != nil {
		return fmt.Errorf("can't fetch xlsx doc sith fetcher: %w", err)
	}
	defer xlsx.Close()

	prevSc, err := s.repo.GetLastSchedule(ctx)
	if err != nil {
		return fmt.Errorf("can't get last schedule from repo: %w", err)
	}

	pr, err := parser.NewParser(xlsx, s.docID)
	if err != nil {
		return fmt.Errorf("can't initialize parser to parse new schedule: %w", err)
	}
	newSc, err := pr.ParseXLSX()
	if err != nil {
		return fmt.Errorf("can't parse xlsx into schedule: %w", err)
	}

	if prevSc == nil {
		err = s.repo.SaveSchedule(ctx, newSc)
		if err != nil {
			return fmt.Errorf("can't save new schedule to repo: %w", err)
		}
		return nil
	}

	// saving to repo parsed schedule
	err = s.repo.SaveSchedule(ctx, newSc)
	if err != nil {
		return fmt.Errorf("can't save new schedule to repo: %w", err)
	}

	update := domain.NewScheduleUpdate(*prevSc, newSc)
	if update.IsEmpty() {
		return nil
	}

	// only publish, if we saved to repo successfully and update is not empty
	s.evBus.Publish(ctx, update)

	return nil
}
