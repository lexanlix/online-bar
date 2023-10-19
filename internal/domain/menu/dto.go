package menu

type CreateMenuDTO struct {
	UserID string             `json:"user_id"`
	Name   string             `json:"name"`
	Drinks map[string][]Drink `json:"drinks"`
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

type AddDrinkDTO struct {
	MenuID     string `json:"menu_id"`
	CategID    uint32 `json:"category_id"`
	IsNewCateg bool   `json:"is_new_category"`
	Drink      Drink  `json:"drink"`
}
