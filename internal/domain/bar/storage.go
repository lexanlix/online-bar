package bar

import "context"

type Repository interface {
	CreateBar(context.Context, CreateBarDTO) (uint32, error)
	CloseBar(context.Context, CloseBarDTO) (string, error)
	UpdateInfo(context.Context, UpdateBarDTO) (string, error)
	GetOrders(context.Context, GetOrdersDTO) ([]string, error)
	GetBarOrders(context.Context, GetBarOrdersDTO) ([]string, error)
}
