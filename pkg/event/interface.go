package event

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrInvalidEventType = errors.New("invalid event type")

type Entity struct {
	ID   uuid.UUID
	Type Type
	Data string
}

type Type string

const (
	TypeCreated Type = "created"
)

//go:generate mockery --name Event
type Event interface {
	SaveEvent(ctx context.Context, entity Entity) error
}

func NewEvent(id uuid.UUID, evType Type, data string) (Entity, error) {
	if err := checkType(evType); err != nil {
		return Entity{}, errors.Wrap(err, "check event type")
	}

	return Entity{
		ID:   id,
		Type: evType,
		Data: data,
	}, nil
}

func checkType(t Type) error {
	if t == TypeCreated {
		return nil
	}

	return ErrInvalidEventType
}
