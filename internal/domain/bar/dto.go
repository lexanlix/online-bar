package bar

type CreateBarDTO struct {
	EventID     string `json:"event_id"`
	Description string `json:"description"`
	Orders      string `json:"orders"`
}
