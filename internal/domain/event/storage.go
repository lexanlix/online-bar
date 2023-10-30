package event

import "context"

type Repository interface {
	CreateEvent(context.Context, CreateEventDTO, []string) (string, error)
	SetActive(context.Context, string) (string, error)
	FindAllUserEvents(context.Context, FindAllEventsDTO) ([]Event, error)
	FindUserEvent(context.Context, FindEventDTO) (Event, error)
	UpdateEvent(context.Context, UpdateEventDTO) (string, error)
	DeleteEvent(context.Context, CompleteEventDTO) (string, error)
}
