package menu

const (
	// Drinks.Category constants
	Beer       = "beers"
	Cider      = "ciders"
	LongDrink  = "long_drinks"
	NonAlco    = "non_alcos"
	ShortDrink = "short_drinks"
	ShotDrink  = "shot_drinks"
	StrongAlco = "strong_alcos"

	// ice_types
	BlockIce   = "block_ice"
	CubedIce   = "cubed_ice"
	CrackedIce = "cracked_ice"
	NuggetIce  = "nugget_ice"
	CrushedIce = "crushed_ice"
	NoIce      = "no_ice"

	// cooking_methods
	Shake = "shake"
	Stir  = "stir"
	Build = "build"
	Blend = "blend"
)

type Menu struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Name      string             `json:"name"`
	Drinks    map[string][]Drink `json:"drinks"`
	TotalCost uint32             `json:"total_cost"`
}

type Drink struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Category       string      `json:"category"`
	Cooking_method string      `json:"cooking_method"`
	Composition    Composition `json:"composition"`
	OrderIceType   string      `json:"order_ice_type"`
	Price          uint32      `json:"price"`
	BarsID         []uint32    `json:"bars_id,omitempty"`
}

// Состав напитка:
// общее количества затраченного льда, ингридиенты разных типов с количеством и размерностью
type Composition struct {
	IceBulk    uint32      `json:"ice_bulk"`
	Liquids    []Liquid    `json:"liquids,omitempty"`
	SolidsBulk []SolidBulk `json:"solids_bulk,omitempty"`
	SolidsUnit []SolidUnit `json:"solids_unit,omitempty"`
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
	Volume uint32 `json:"volume"`
}
