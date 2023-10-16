package menu

import "context"

type Repository interface {
	CreateMenu(context.Context, CreateMenuDTO) (string, error)
	DeleteMenu(context.Context, DeleteMenuDTO) (string, error)
	FindMenu(context.Context, FindMenuDTO) (Menu, error)
	UpdateMenu(context.Context, UpdateMenuDTO) (string, error)
}
