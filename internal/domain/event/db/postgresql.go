package db

import (
	"context"
	"restapi/internal/domain/event"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

// Create implements event.Repository.
func (r *repository) Create(ctx context.Context, event *event.Event) error {
	panic("unimplemented")
}

// Delete implements event.Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements event.Repository.
func (r *repository) FindAll(ctx context.Context) ([]event.Event, error) {
	panic("unimplemented")
}

// FindOne implements event.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (event.Event, error) {
	panic("unimplemented")
}

// Update implements event.Repository.
func (r *repository) Update(ctx context.Context, event event.Event) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) event.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
