package event

import (
	"restapi/internal/domain/report"
	"time"
)

type Event struct {
	ID                 string        `json:"id"`
	HostID             string        `json:"host_id"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	ParticipantsNumber uint32        `json:"participants_number"`
	DateTime           time.Time     `json:"date_time"`
	Status             string        `json:"status"`
	MenuID             string        `json:"menu_id"`
	ShoppingList       []string      `json:"shopping_list"`
	Report             report.Report `json:"report"`
}
