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

type ScheduleService struct {
	fetcher XLSXScheduleDocumentFetcher
	docID   string
}

func NewScheduleService(fetcher XLSXScheduleDocumentFetcher, docID string) *ScheduleService {
	return &ScheduleService{
		fetcher: fetcher,
	}
}

func (s *ScheduleService) GetSchedule(ctx context.Context) (domain.Schedule, error) {
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
