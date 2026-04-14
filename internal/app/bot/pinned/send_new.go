package pinned

import (
	"context"
	"fmt"
	"time"

	"github.com/hurtki/school-events-bot/internal/domain"
	bot_domain "github.com/hurtki/school-events-bot/internal/domain/bot"
)

// createNewPinnedMessage handles the sending, pinning, and state storage of a new message.
func (s *BotUpcomingEventsPinService) createNewPinnedMessage(ctx context.Context, summary domain.UpcomingEventsSummary) error {
	msgID, err := s.bot.SendEventsSummaryAndPin(summary)
	if err != nil {
		return fmt.Errorf("bot failed to send and pin summary: %w", err)
	}

	pinnedMsg, err := bot_domain.NewPinnedMessage(msgID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to initialize pinned message entity: %w", err)
	}

	if err := s.pinnedMsgRepo.SaveLastPinnedMessage(ctx, pinnedMsg); err != nil {
		return fmt.Errorf("failed to persist pinned message state: %w", err)
	}

	s.logger.Info("successfully created and registered new pinned message", "msg_id", msgID)
	return nil
}
