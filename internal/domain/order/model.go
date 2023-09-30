package order

import "time"

type Order struct {
	ID        string    `json:"id"`
	BarID     string    `json:"bar_id"`
	OrderBody OrderBody `json:"order_body"`
	DateTime  time.Time `json:"date_time"`
}

type OrderBody struct {
	DrinksID []string
	Comment  string
}
