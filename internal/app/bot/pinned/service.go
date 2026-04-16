package pinned

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/domain"
	bot_domain "github.com/hurtki/school-events-bot/internal/domain/bot"
)

const (
	// upcomingEventsShowCount defines how many events to include in the summary for every group
	upcomingEventsShowCount = 3
	// maxPinnedMessageAge defines the threshold after which a message is considered stale and should be re-sent
	maxPinnedMessageAge = 48 * time.Hour
	// maxPinnedMessageAge = time.Second * 90
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
	bot           *bot.Bot
}

func NewBotUpcomingEventsPinService(
	logger *slog.Logger,
	pinnedMsgRepo PinnedMessageStateRepo,
	scheduleRepo ScheduleRepository,
	bot *bot.Bot,
) *BotUpcomingEventsPinService {
	return &BotUpcomingEventsPinService{
		logger:        logger.With("service", "bot-upcoming-events-service"),
		pinnedMsgRepo: pinnedMsgRepo,
		scheduleRepo:  scheduleRepo,
		bot:           bot,
	}
}

// Update orchestrates the pinned message lifecycle: fetching data, deciding whether to
// edit or re-send, and persisting the state.
func (s *BotUpcomingEventsPinService) Update(ctx context.Context) error {
	// 1. Fetch latest schedule data
	schedule, err := s.scheduleRepo.GetLastSchedule(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last schedule from repo: %w", err)
	} else if schedule == nil {
		return fmt.Errorf("there is no last schedule in repo")
	}

	// 2. Generate content summary based on specific event types
	summary := schedule.GetUpcomingEventsSummary(
		upcomingEventsShowCount,
		domain.BagrutTestEvent,
		domain.ProtectionBagrutTestEvent,
	)

	// 3. Retrieve current state of the pinned message
	msgInfo, exists, err := s.pinnedMsgRepo.GetLastPinnedMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pinned message state: %w", err)
	}

	// 4. Determine if we need to create a brand new message (if none exists or if it's too old)
	shouldCreateNew := !exists || time.Since(msgInfo.LastMessageSentAt) > maxPinnedMessageAge

	if shouldCreateNew {
		if exists {
			// Clean up the old message if it exists. Error is ignored as it might be already deleted.
			err = s.bot.DeleteMessage(msgInfo.LastMessageID)
			if err != nil {
				s.logger.Warn("can't delete message", "err", err)
			}
		}
		return s.createNewPinnedMessage(ctx, summary)
	}

	// 5. Attempt to update the existing message
	err = s.bot.UpdateMessageWithEventsSummary(msgInfo.LastMessageID, summary)
	if err != nil {
		s.logger.Warn("failed to update existing message, falling back to recreation", "err", err)
		return s.createNewPinnedMessage(ctx, summary)
	}

	return nil
}

// HandleScheduleUpdate acts as a subscriber/callback for schedule changes.
func (s *BotUpcomingEventsPinService) HandleScheduleUpdate(ctx context.Context, update domain.ScheduleUpdate) {
	if update.IsEmpty() {
		s.logger.Debug("received empty schedule update, skipping")
		return
	}
	s.logger.Info("received update, handling", "new-events-count", len(update.Added), "delete-events-count", len(update.Deleted))

	if err := s.Update(ctx); err != nil {
		s.logger.Error("failed to process schedule update", "err", err)
	}
}
