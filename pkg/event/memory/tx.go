package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/veleton777/booking_api/pkg/event"
)

type EventClient struct {
	storage map[event.Type]map[uuid.UUID]event.Entity

	mu *sync.Mutex
}

func NewEventClient() *EventClient {
	return &EventClient{
		mu:      &sync.Mutex{},
		storage: make(map[event.Type]map[uuid.UUID]event.Entity),
	}
}

func (e *EventClient) SaveEvent(_ context.Context, entity event.Entity) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.storage[entity.Type]) == 0 {
		e.storage[entity.Type] = make(map[uuid.UUID]event.Entity)
	}

	e.storage[entity.Type][entity.ID] = entity

	return nil
}

func (e *EventClient) AckEvent(_ context.Context, entity event.Entity) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.storage[entity.Type], entity.ID)

	return nil
}

func (e *EventClient) GetEventsByType(_ context.Context, evType event.Type, limit int) ([]event.Entity, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if limit > len(e.storage[evType]) {
		limit = len(e.storage[evType])
	}

	res := make([]event.Entity, 0, limit)
	for _, v := range e.storage[evType] {
		res = append(res, v)
	}

	return res, nil
}
