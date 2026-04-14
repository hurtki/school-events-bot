package updates

import (
	"context"
	"log/slog"

	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/domain"
)

type BotScheduleUpdatesService struct {
	logger *slog.Logger
	bot    *bot.Bot
}

func NewBotScheduleUpdatesService(logger *slog.Logger, bot *bot.Bot) *BotScheduleUpdatesService {
	return &BotScheduleUpdatesService{
		bot:    bot,
		logger: logger,
	}
}

func (s *BotScheduleUpdatesService) HandleScheduleUpdate(ctx context.Context, update domain.ScheduleUpdate) {
	s.logger.Debug("handling schedule update")
	err := s.bot.NotifyAboutUpdate(ctx, update)
	if err != nil {
		s.logger.Error("error occured, when notifying using bot", "err", err)
		return
	}
	s.logger.Debug("notified about the update successfully")
}
