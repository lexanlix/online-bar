package bar

type Bar struct {
	ID          string   `json:"id"`
	EventID     string   `json:"event_id"`
	Name        string   `json:"name"`
	Description string   `json:"info"`
	Orders      []string `json:"orders"` // массив id-ов заказов
	Status      string   `json:"status"`
	SessionURL  string
}
