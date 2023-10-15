package bar

type CreateBarDTO struct {
	EventID     string `json:"event_id"`
	Name        string `json:"name"`
	Description string `json:"info"`
}

type RespCreateBar struct {
	ID uint32 `json:"id"`
}

type CloseBarDTO struct {
	ID uint32 `json:"id"`
}

type UpdateBarDTO struct {
	ID          uint32   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"info"`
	Orders      []string `json:"orders"`
	SessionURL  string   `json:"session_url"`
}

type GetOrdersDTO struct {
	EventID string `json:"event_id"`
}

type GetBarOrdersDTO struct {
	ID      uint32 `json:"id"`
	EventID string `json:"event_id"`
}
