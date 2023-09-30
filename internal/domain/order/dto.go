package order

type CreateOrderDTO struct {
	BarID     string    `json:"bar_id"`
	OrderBody OrderBody `json:"order_body"`
}
