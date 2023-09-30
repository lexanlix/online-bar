package event

import "context"

type Repository interface {
	Create(ctx context.Context, event *Event) error
	FindAll(ctx context.Context) ([]Event, error)
	FindOne(ctx context.Context, id string) (Event, error)
	Update(ctx context.Context, event Event) error
	Delete(ctx context.Context, id string) error
}
