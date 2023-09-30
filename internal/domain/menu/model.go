package menu

// Drinks.Group constants
const (
	Cocktails  = "cocktail"
	Shots      = "shot"
	NonAlco    = "nonAlco"
	Beer       = "beer"
	StrongAlco = "strongAlco"
)

type Menu interface {
	AddDrink(newDrink Drink) error
	UpdateTotalCost()
}

type menu struct {
	Drinks    []Drink `json:"drinks"`
	TotalCost uint32  `json:"total_cost"`
}

type Drink struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Price uint32 `json:"price"`
}

func NewMenu(drinks []Drink) Menu {
	return &menu{
		Drinks: drinks,
	}
}

func (m *menu) AddDrink(newDrink Drink) error {
	m.Drinks = append(m.Drinks, newDrink)
	m.UpdateTotalCost()

	return nil
}

func (m *menu) UpdateTotalCost() {
	for _, drink := range m.Drinks {
		m.TotalCost += drink.Price
	}
}
