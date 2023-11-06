package event

import "context"

type Repository interface {
	CreateEvent(context.Context, CreateEventDTO, []string) (string, error)
	SetActive(context.Context, string) error
	FindAllUserEvents(context.Context, FindAllEventsDTO) (RespAllEvents, error)
	FindUserEvent(context.Context, FindEventDTO) (Event, error)
	UpdateEvent(context.Context, UpdateEventDTO) error
	DeleteEvent(context.Context, CompleteEventDTO) error
	UpdateIceTypesNum(context.Context, bool, string) error
	GetIceTypesNum(context.Context, string) (bool, error)
}
