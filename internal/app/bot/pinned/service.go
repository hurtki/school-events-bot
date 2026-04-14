package pinned

import (
	"context"
	"log/slog"

	"github.com/hurtki/school-events-bot/internal/domain"
	bot_domain "github.com/hurtki/school-events-bot/internal/domain/bot"
)

type PinnedMessageStateRepo interface {
	GetLastPinnedMessage(ctx context.Context) (bot_domain.PinnedMessage, bool, error)
	SaveLastPinnedMessage(ctx context.Context, msg bot_domain.PinnedMessage) error
}

type ScheduleRepository interface {
	GetLastSchedule(ctx context.Context) (*domain.Schedule, error)
}

type BotUpcomingEventsPinService struct {
	logger        *slog.Logger
	pinnedMsgRepo PinnedMessageStateRepo
	scheduleRepo  ScheduleRepository
}

func NewBotUpcomingEventsPinService(
	logger *slog.Logger,
	pinnedMsgRepo PinnedMessageStateRepo,
	scheduleRepo ScheduleRepository,
) *BotUpcomingEventsPinService {
	return &BotUpcomingEventsPinService{
		logger:        logger,
		pinnedMsgRepo: pinnedMsgRepo,
		scheduleRepo:  scheduleRepo,
	}
}

func (s *BotUpcomingEventsPinService) HandleScheduleUpdate(ctx context.Context, update domain.ScheduleUpdate) {
	// subscription on update events
}

func (s *BotUpcomingEventsPinService) Update(ctx context.Context) error {
	// is called by worker
	return nil
}
