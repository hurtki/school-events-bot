package evbus

import (
	"context"

	evBus "github.com/asaskevich/EventBus"
	"github.com/hurtki/school-events-bot/internal/domain"
)

type ScheduleUpdateEventBus struct {
	bus       evBus.Bus
	topicName string
}

type HandleEventFunc func(ctx context.Context, event domain.ScheduleUpdate)

func NewScheduleUpdateEventBus() *ScheduleUpdateEventBus {
	bus := evBus.New()
	return &ScheduleUpdateEventBus{
		bus:       bus,
		topicName: "topic:scheduleUpdates",
	}
}

func (b *ScheduleUpdateEventBus) Publish(ctx context.Context, event domain.ScheduleUpdate) {
	b.bus.Publish(b.topicName, ctx, event)
}

func (b *ScheduleUpdateEventBus) Subscribe(f HandleEventFunc) {
	b.bus.SubscribeAsync(b.topicName, f, false)
}
