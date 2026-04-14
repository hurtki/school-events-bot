package evbus

import (
	"context"
	"sync"
	"sync/atomic"

	evBus "github.com/asaskevich/EventBus"
	"github.com/hurtki/school-events-bot/internal/domain"
)

// ScheduleUpdateEventBus wraps event bus to wait for all events to be proceeded
// Subscribe shouldn't be called after Publish!!!
type ScheduleUpdateEventBus struct {
	bus       evBus.Bus
	topicName string

	subscribersCount atomic.Int32
	mu               sync.RWMutex
	i                atomic.Int64
}

type HandleEventFunc func(ctx context.Context, update domain.ScheduleUpdate)

func NewScheduleUpdateEventBus() *ScheduleUpdateEventBus {
	bus := evBus.New()
	return &ScheduleUpdateEventBus{
		bus:       bus,
		topicName: "topic:scheduleUpdates",
	}
}

func (b *ScheduleUpdateEventBus) Publish(ctx context.Context, update domain.ScheduleUpdate) {
	wg := &sync.WaitGroup{}

	wg.Add(int(b.subscribersCount.Load()))
	b.bus.Publish(b.topicName, ctx, update, wg)

	wg.Wait()
}

func (b *ScheduleUpdateEventBus) Subscribe(f HandleEventFunc) {
	b.subscribersCount.Add(1)
	b.bus.SubscribeAsync(b.topicName, func(ctx context.Context, update domain.ScheduleUpdate, wg *sync.WaitGroup) {
		defer wg.Done()
		f(ctx, update)
	}, false)
}
