package menu

import "context"

type Repository interface {
	CreateMenu(context.Context, CreateMenuDTO, uint32) (string, error)
	DeleteMenu(context.Context, DeleteMenuDTO) error
	FindMenu(context.Context, FindMenuDTO) (Menu, error)
	UpdateMenu(context.Context, UpdateMenuDTO, uint32) (string, error)
	UpdateNameMenu(context.Context, UpdateMenuNameDTO) error
	AddDrink(context.Context, AddDrinkDTO) error
}
