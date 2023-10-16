package menu

type CreateMenuDTO struct {
	UserID string             `json:"user_id"`
	Name   string             `json:"name"`
	Drinks map[string][]Drink `json:"drinks"`
}

type DeleteMenuDTO struct {
	ID string `json:"id"`
}

type FindMenuDTO struct {
	ID string `json:"id"`
}

type UpdateMenuDTO struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Name      string             `json:"name"`
	Drinks    map[string][]Drink `json:"drinks"`
	TotalCost uint32             `json:"total_cost"`
}

type GetMapDTO struct {
	Category string  `json:"category"`
	Drinks   []Drink `json:"drinks"`
}

type AddDrinkDTO struct {
	ID             uint32      `json:"id"`
	Name           string      `json:"name"`
	Category       string      `json:"category"`
	Cooking_method string      `json:"cooking_method"`
	Composition    Composition `json:"composition"`
	IceType        string      `json:"ice_type"`
	Price          uint32      `json:"price"`
	BarsID         []uint32    `json:"bars_id"`
}
