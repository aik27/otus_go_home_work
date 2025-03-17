package storage

import (
	"context"
)

type Repository interface {
	Save(event Event) (Event, error)
	Delete(event Event) error
	GetByUserId(userId int64) (EventList, error)
	GetAll() (EventList, error)
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}
