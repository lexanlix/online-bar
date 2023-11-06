package bar

import "context"

type Repository interface {
	CreateBar(context.Context, CreateBarDTO) (uint32, error)
	CloseBar(context.Context, CloseBarDTO) error
	UpdateInfo(context.Context, UpdateBarDTO) error
	GetOrders(context.Context, GetOrdersDTO) ([]string, error)
	GetBarOrders(context.Context, GetBarOrdersDTO) ([]string, error)
}
