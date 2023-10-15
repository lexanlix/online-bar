package event

import "context"

type Repository interface {
	CreateEvent(context.Context, CreateEventDTO) (string, error)
	SetActive(context.Context, string) (string, error)
	FindAllUserEvents(context.Context, FindAllEventsDTO) ([]Event, error)
	FindOneUserEvent(context.Context, FindEventDTO) (Event, error)
	UpdateEvent(context.Context, UpdateEventDTO) (string, error)
	DeleteEvent(context.Context, CompleteEventDTO) (string, error)
}
