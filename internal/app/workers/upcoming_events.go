package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/hurtki/school-events-bot/internal/app/bot/pinned"
	"github.com/hurtki/school-events-bot/internal/bot"
)

type UpcomingEventsWorker struct {
	service  *pinned.BotUpcomingEventsPinService
	bot      *bot.Bot
	interval time.Duration
	logger   *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewUpcomingEventsWorker(logger *slog.Logger, service *pinned.BotUpcomingEventsPinService, interval time.Duration) *UpcomingEventsWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &UpcomingEventsWorker{
		service:  service,
		interval: interval,
		logger:   logger.With("service", "upcoming-events-worker"),
		ctx:      ctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
	}
}

func (w *UpcomingEventsWorker) Close(ctx context.Context) error {
	w.cancel()
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		w.logger.Warn("couldn't shutdown in time, exiting", "ctxErr", ctx.Err())
		return ctx.Err()
	case <-done:
		w.logger.Info("successfully shutted down")
		return nil
	}
}

func (w *UpcomingEventsWorker) Start() {
	w.wg.Go(w.run)
}

func (w *UpcomingEventsWorker) run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.logger.Info("started polling", "interval", w.interval.String())
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(w.ctx, time.Second*10)
			defer cancel()

			w.logger.Info("poll unit started")
			start := time.Now()
			err := w.service.Update(ctx)
			if err != nil {
				w.logger.Error("error occured", "err", err, "duration", time.Since(start).String())
			} else {
				w.logger.Info("updated upcoming events", "duration", time.Since(start).String())
			}
		}
	}
}
