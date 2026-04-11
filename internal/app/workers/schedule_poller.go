package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/hurtki/school-events-bot/internal/app/schedule"
	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/domain"
)

type SchedulePoller struct {
	service  *schedule.ScheduleService
	bot      *bot.Bot
	interval time.Duration
	logger   *slog.Logger
	repo     ScheduleRepository

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type ScheduleRepository interface {
	GetLastSchedule(ctx context.Context) (*domain.Schedule, error)
	SaveSchedule(ctx context.Context, s domain.Schedule) error
}

func NewSchedulePoller(logger *slog.Logger, service *schedule.ScheduleService, bot *bot.Bot, interval time.Duration, repo ScheduleRepository) *SchedulePoller {
	ctx, cancel := context.WithCancel(context.Background())
	return &SchedulePoller{
		service:  service,
		bot:      bot,
		interval: interval,
		logger:   logger.With("service", "schedule-poller"),
		ctx:      ctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
		repo:     repo,
	}
}

func (p *SchedulePoller) Close(ctx context.Context) error {
	p.cancel()
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		p.logger.Warn("couldn't shutdown in time, exiting", "ctxErr", ctx.Err())
		return ctx.Err()
	case <-done:
		p.logger.Info("successfully shutted down")
		return nil
	}
}

func (p *SchedulePoller) Start() {
	p.wg.Go(p.run)
}

func (p *SchedulePoller) run() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	p.logger.Info("started polling", "interval", p.interval.String())
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(p.ctx, time.Second*10)
			defer cancel()

			p.logger.Info("poll unit started")
			sc, err := p.service.GetSchedule(ctx)
			if err != nil {
				p.logger.Error("can't get schedule from service", "err", err)
				continue
			}
			prevSc, err := p.repo.GetLastSchedule(ctx)
			if err != nil {
				p.logger.Error("can't get previous schedule", "err", err)
			}
			if prevSc == nil {
				p.logger.Info("previous schedule wasn't saved, writing new")
				err = p.repo.SaveSchedule(ctx, sc)
				if err != nil {
					p.logger.Error("can't save schedule", "err", err)
				}
				continue
			}

			// if we got prevoius and new schedules successfully
			// we will compare them
			update := domain.CompareSchedules(*prevSc, sc)

			if update.IsEmpty() {
				p.logger.Info("nothing changed since last update, not notifying about updates")
			} else {
				p.logger.Info("there are new updates", "deleted", len(update.Deleted), "added", len(update.Added))
				err = p.bot.NotifyAboutUpdate(ctx, update)
				if err != nil {
					p.logger.Error("can't notify about update", "err", err)
				} else {
					p.logger.Info("notified about new update successfully")
				}
				err = p.repo.SaveSchedule(ctx, sc)
				if err != nil {
					p.logger.Error("can't save schedule", "err", err)
				}
			}
		}
	}
}
