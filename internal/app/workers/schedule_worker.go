package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/hurtki/school-events-bot/internal/app/schedule"
	"github.com/hurtki/school-events-bot/internal/bot"
)

type ScheduleWorker struct {
	service  *schedule.ScheduleService
	bot      *bot.Bot
	interval time.Duration
	logger   *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewScheduleWorker(logger *slog.Logger, service *schedule.ScheduleService, interval time.Duration) *ScheduleWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &ScheduleWorker{
		service:  service,
		interval: interval,
		logger:   logger.With("service", "schedule-worker"),
		ctx:      ctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
	}
}

func (p *ScheduleWorker) Close(ctx context.Context) error {
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

func (p *ScheduleWorker) Start() {
	p.wg.Go(p.run)
}

func (p *ScheduleWorker) run() {
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
			start := time.Now()
			err := p.service.Update(ctx)
			if err != nil {
				p.logger.Error("error occured", "err", err, "duration", time.Since(start).String())
			} else {
				p.logger.Info("updated shedule service", "duration", time.Since(start).String())
			}
		}
	}
}

// if err != nil {
// 	p.logger.Error("can't get schedule from service", "err", err)
// 	continue
// }
// if err != nil {
// 	p.logger.Error("can't get previous schedule", "err", err)
// }
// if prevSc == nil {
// 	p.logger.Info("previous schedule wasn't saved, writing new")
// 	err = p.repo.SaveSchedule(ctx, sc)
// 	if err != nil {
// 		p.logger.Error("can't save schedule", "err", err)
// 	}
// 	continue
// }
//
// // if we got prevoius and new schedules successfully
// // we will compare them
// update := domain.CompareSchedules(*prevSc, sc)
//
// if update.IsEmpty() {
// 	p.logger.Info("nothing changed since last update, not notifying about updates")
// } else {
// 	p.logger.Info("there are new updates", "deleted", len(update.Deleted), "added", len(update.Added))
// 	err = p.bot.NotifyAboutUpdate(ctx, update)
// 	if err != nil {
// 		p.logger.Error("can't notify about update", "err", err)
// 	} else {
// 		p.logger.Info("notified about new update successfully")
// 	}
// 	err = p.repo.SaveSchedule(ctx, sc)
// 	if err != nil {
// 		p.logger.Error("can't save schedule", "err", err)
// 	}
// }
