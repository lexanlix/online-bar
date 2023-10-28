package drinks_list

import (
	"context"
	"restapi/internal/domain/menu"
)

type Repository interface {
	AddUserDrink(context.Context, AddUserDrinkDTO) (string, error)
	DeleteUserDrink(context.Context, DeleteUserDrinkDTO) error
	FindUserDrink(context.Context, FindUserDrinkDTO) (menu.Drink, error)
	FindUserDrinks(context.Context, FindUserDrinksDTO) (RespFindUDrinks, error)
	UpdateUserDrink(context.Context, UpdateUserDrinkDTO) (string, error)
}
