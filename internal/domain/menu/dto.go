package menu

type MenuDTO struct {
	UserID string                   `json:"user_id"`
	Name   string                   `json:"name"`
	Drinks map[string][]NewDrinkDTO `json:"drinks,omitempty"`
}

type CreateMenuDTO struct {
	UserID    string   `json:"user_id"`
	Name      string   `json:"name"`
	DrinksIDs []string `json:"drinks_id,omitempty"`
}

type NewDrinkDTO struct {
	Name           string      `json:"name"`
	Category       string      `json:"category"`
	Cooking_method string      `json:"cooking_method"`
	Composition    Composition `json:"composition"`
	OrderIceType   string      `json:"order_ice_type"`
	Price          uint32      `json:"price"`
	BarsID         []uint32    `json:"bars_id,omitempty"`
}

type UserMenusDTO struct {
	UserID string `json:"user_id"`
}

type RespCreateMenuDTO struct {
	ID string `json:"id"`
}

type DeleteMenuDTO struct {
	ID string `json:"id"`
}

type FindMenuDTO struct {
	ID string `json:"id"`
}

type UpdateMenuDTO struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Drinks map[string][]Drink `json:"drinks"`
}

type UpdateMenuNameDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RespUserMenus struct {
	Menus []UserMenu `json:"menus"`
}

type UserMenu struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AddDrinkDTO struct {
	MenuID string      `json:"menu_id"`
	Drink  NewDrinkDTO `json:"drink"`
}

type AddDrinkFromListDTO struct {
	MenuID  string `json:"menu_id"`
	DrinkID string `json:"drink_id"`
}

type DeleteDrinkDTO struct {
	DrinkID string `json:"drink_id"`
}
