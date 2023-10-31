package event

import "time"

type CreateEventDTO struct {
	UserID             string    `json:"user_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ParticipantsNumber uint32    `json:"participants_number"`
	DateTime           time.Time `json:"date_time"`
	MenuID             string    `json:"menu_id"`
}

type CompleteEventDTO struct {
	ID string `json:"id"`
}

type FindAllEventsDTO struct {
	UserID string `json:"user_id"`
}

type RespAllEvents struct {
	Events []Event `json:"events"`
}

type FindEventDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type UpdateEventDTO struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ParticipantsNumber uint32    `json:"participants_number"`
	DateTime           time.Time `json:"date_time"`
}

type RespCreateEvent struct {
	ID string `json:"id"`
}
