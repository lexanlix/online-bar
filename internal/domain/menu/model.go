package menu

const (
	// Drinks.Category constants
	ShotDrink  = "shot_drinks"
	ShortDrink = "short_drinks"
	LongDrink  = "long_drinks"
	NonAlco    = "non_alcos"
	Beer       = "beers"
	Cider      = "ciders"
	StrongAlco = "strong_alcos"

	// ice_types
	CrashIce = "crash_ice"
	CubeIce  = "cube_ice"
	NoIce    = "no_ice"

	// cooking_methods
	Shake = "shake"
	Stir  = "stir"
	Build = "build"
	Blend = "blend"
)

// type Menu interface {
// 	AddDrink(dto AddDrinkDTO) error
// 	UpdateTotalCost()
// }

type Menu struct {
	ID        uint32             `json:"id"`
	Name      string             `json:"name"`
	Drinks    map[string][]Drink `json:"drinks"`
	TotalCost uint32             `json:"total_cost"`
}

type Drink struct {
	ID             uint32      `json:"id"`
	Name           string      `json:"name"`
	Category       string      `json:"category"`
	Cooking_method string      `json:"cooking_method"`
	Composition    Composition `json:"composition"`
	IceType        string      `json:"ice_type"`
	Price          uint32      `json:"price"`
	BarsID         []uint32    `json:"bars_id"`
}

// TODO состав напитка
type Composition struct {
}

func NewMenu(id uint32, name string, drinks map[string][]Drink) Menu {
	return Menu{
		ID:     id,
		Name:   name,
		Drinks: drinks,
	}
}

func (m *Menu) AddDrink(dto Drink) error {
	newDrink := Drink(dto)

	m.Drinks[dto.Category] = append(m.Drinks[dto.Category], newDrink)
	m.UpdateTotalCost()

	return nil
}

func (m *Menu) UpdateTotalCost() {
	for category := range m.Drinks {
		for _, drink := range m.Drinks[category] {
			m.TotalCost += drink.Price
		}
	}
}
