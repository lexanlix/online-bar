package bar

import (
	"restapi/internal/domain/menu"
)

type Bar struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	Description string    `json:"info"`
	Orders      []string  `json:"orders"` // массив id-ов заказов
	Menu        menu.Menu `json:"menu"`
	Status      string    `json:"status"`
}
