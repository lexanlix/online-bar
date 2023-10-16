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
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
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
	OrderIceType   string      `json:"order_ice_type"`
	Price          uint32      `json:"price"`
	BarsID         []uint32    `json:"bars_id"`
}

// Состав напитка:
// общее количества затраченного льда, ингридиенты разных типов с количеством и размерностью
type Composition struct {
	IceBulk    uint32      `json:"ice_bulk"`
	Liquids    []Liquid    `json:"liquids"`
	SolidsBulk []SolidBulk `json:"solids_bulk"`
	SolidsUnit []SolidUnit `json:"solids_unit"`
}

// Жидкие ингридиенты, имеющие объем и его ед. изм.
type Liquid struct {
	Name   string `json:"name"`
	Unit   string `json:"unit"`
	Volume uint32 `json:"volume"`
}

// Твердые ингридиенты, имеющие объем и его ед. изм.
type SolidBulk struct {
	Name   string `json:"name"`
	Unit   string `json:"unit"`
	Volume uint32 `json:"volume"`
}

// Твердые ингридиенты, считающиеся в штуках
type SolidUnit struct {
	Name   string `json:"name"`
	Amount uint32 `json:"amount"`
}

func NewMenu(id string, name string, drinks map[string][]Drink) Menu {
	return Menu{
		ID:     id,
		Name:   name,
		Drinks: drinks,
	}
}

// func (m *Menu) AddDrink(dto Drink) error {
// 	newDrink := Drink(dto)

// 	m.Drinks[dto.Category] = append(m.Drinks[dto.Category], newDrink)
// 	m.UpdateTotalCost()

// 	return nil
// }

func UpdateTotalCost(m *Menu) {
	for category := range m.Drinks {
		for _, drink := range m.Drinks[category] {
			m.TotalCost += drink.Price
		}
	}
}

// TODO
func GetTotalCost(m *UpdateMenuDTO) uint32 {
	for category := range m.Drinks {
		for _, drink := range m.Drinks[category] {
			m.TotalCost += drink.Price
		}
	}
	return 0
}
