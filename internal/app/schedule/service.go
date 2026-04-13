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
	// updates information in repository
	// and sends events about updates of schedule to event bus
	// if there are

	return nil
}

func (s *ScheduleService) getSchedule(ctx context.Context) (domain.Schedule, error) {
	xlsx, err := s.fetcher.FetchXLSX(ctx, s.docID)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("can't fetch xlsx: %w", err)
	}
	defer func() {
		err := xlsx.Close()
		if err != nil {
			fmt.Println("can't close xlsx")
		}
	}()
	p, err := parser.NewParser(xlsx, s.docID)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("can't parse xlsx: %w", err)
	}
	return p.ParseXLSX()
}
