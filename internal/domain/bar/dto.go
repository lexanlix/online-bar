package bar

import "restapi/internal/domain/menu"

type CreateBarDTO struct {
	EventID       string             `json:"event_id"`
	Description   string             `json:"info"`
	CreateMenuDTO menu.CreateMenuDTO `json:"menu"`
}

type RespCreateBar struct {
	ID   uint32    `json:"id"`
	Menu menu.Menu `json:"menu"`
}

type CloseBarDTO struct {
	ID uint32 `json:"id"`
}

type UpdateBarDTO struct {
	ID          uint32   `json:"id"`
	Description string   `json:"info"`
	Orders      []string `json:"orders"`
}

type GetOrdersDTO struct {
	EventID string `json:"event_id"`
}

type GetBarOrdersDTO struct {
	ID      uint32 `json:"id"`
	EventID string `json:"event_id"`
}
