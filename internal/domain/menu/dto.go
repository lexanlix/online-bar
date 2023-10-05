package menu

type CreateMenuDTO struct {
	Drinks map[string][]Drink `json:"drinks"`
}

type AddDrinkDTO struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Price uint32 `json:"price"`
}
