package drinks_list

import "restapi/internal/domain/menu"

type AddUserDrinkDTO struct {
	UserID         string           `json:"user_id"`
	Name           string           `json:"name"`
	Category       string           `json:"category"`
	Cooking_method string           `json:"cooking_method"`
	Composition    menu.Composition `json:"composition"`
	OrderIceType   string           `json:"order_ice_type"`
	Price          uint32           `json:"price"`
	BarsID         []uint32         `json:"bars_id,omitempty"`
}

type DeleteUserDrinkDTO struct {
	ID string `json:"id"`
}

type FindUserDrinkDTO struct {
	ID string `json:"id"`
}

type FindUserDrinksDTO struct {
	UserID string `json:"user_id"`
}

type RespFindUDrinks struct {
	Drinks []menu.Drink `json:"drinks"`
}

type UpdateUserDrinkDTO struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Category       string           `json:"category"`
	Cooking_method string           `json:"cooking_method"`
	Composition    menu.Composition `json:"composition"`
	OrderIceType   string           `json:"order_ice_type"`
	Price          uint32           `json:"price"`
	BarsID         []uint32         `json:"bars_id,omitempty"`
}
