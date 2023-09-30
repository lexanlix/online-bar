package bar

import (
	"restapi/internal/domain/menu"
	"restapi/internal/domain/order"
)

type Bar struct {
	ID          string        `json:"id"`
	EventID     string        `json:"event_id"`
	Description string        `json:"description"`
	Orders      []order.Order `json:"orders"`
	Menu        menu.Menu     `json:"-"`
}
