package menu

import "context"

type Repository interface {
	CreateMenu(context.Context, MenuDTO, uint32) (string, error)
	FindUserDrinks(context.Context, []string) ([]Drink, error)
	DeleteMenu(context.Context, DeleteMenuDTO) error
	FindMenu(context.Context, FindMenuDTO) (Menu, error)
	FindUserMenus(context.Context, UserMenusDTO) (RespUserMenus, error)
	UpdateMenu(context.Context, UpdateMenuDTO, uint32) (string, error)
	UpdateNameMenu(context.Context, UpdateMenuNameDTO) error
	AddDrink(context.Context, AddDrinkDTO) (string, error)
	DeleteDrink(context.Context, DeleteDrinkDTO) error
	FindUserDrink(context.Context, string) (NewDrinkDTO, error)
}
