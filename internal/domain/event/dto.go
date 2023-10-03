package event

import "time"

type CreateEventDTO struct {
	HostID             string    `json:"host_id"`
	Name               string    `json:"name"`
	Description        string    `json:"info"`
	ParticipantsNumber uint32    `json:"participants_number"`
	DateTime           time.Time `json:"date_time"`
}

type DeleteEventDTO struct {
	ID string `json:"id"`
}

type FindAllEventsDTO struct {
	HostID string `json:"host_id"`
}

type FindEventDTO struct {
	ID     string `json:"id"`
	HostID string `json:"host_id"`
}

type UpdateEventDTO struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"info"`
	ParticipantsNumber uint32    `json:"participants_number"`
	DateTime           time.Time `json:"date_time"`
}

type RespCreateEvent struct {
	ID string `json:"id"`
}
