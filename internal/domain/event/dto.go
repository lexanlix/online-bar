package event

import "time"

type CreateEventDTO struct {
	HostID             string    `json:"host_id"`
	Name               string    `json:"name"`
	Description        string    `json:"info"`
	ParticipantsNumber uint32    `json:"participants_number"`
	DateTime           time.Time `json:"date_time"`
}
