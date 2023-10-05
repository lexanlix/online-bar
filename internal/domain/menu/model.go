package menu

// Drinks.Group constants
const (
	Cocktails  = "cocktails"
	Shots      = "shots"
	NonAlco    = "nonAlcos"
	Beer       = "beers"
	StrongAlco = "strongAlcos"
)

// type Menu interface {
// 	AddDrink(dto AddDrinkDTO) error
// 	UpdateTotalCost()
// }

type Menu struct {
	Drinks    map[string][]Drink `json:"drinks"`
	TotalCost uint32             `json:"total_cost"`
}

type Drink struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func NewMenu(drinks map[string][]Drink) Menu {
	return Menu{
		Drinks: drinks,
	}
}

func (m *Menu) AddDrink(dto AddDrinkDTO) error {
	newDrink := Drink{
		ID:    dto.ID,
		Name:  dto.Name,
		Price: dto.Price,
	}

	m.Drinks[dto.Group] = append(m.Drinks[dto.Group], newDrink)
	m.UpdateTotalCost()

	return nil
}

func (m *Menu) UpdateTotalCost() {
	for group := range m.Drinks {
		for _, drink := range m.Drinks[group] {
			m.TotalCost += drink.Price
		}
	}
}
