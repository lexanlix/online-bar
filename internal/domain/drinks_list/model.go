package drinks_list

import "restapi/internal/domain/menu"

const (
	// Drinks.Category constants
	Beer       = "beers"
	Cider      = "ciders"
	LongDrink  = "long_drinks"
	NonAlco    = "non_alcos"
	ShortDrink = "short_drinks"
	ShotDrink  = "shot_drinks"
	StrongAlco = "strong_alcos"

	// order_ice_types
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

type UserDrink struct {
	ID             string           `json:"id"`
	UserID         string           `json:"user_id"`
	Name           string           `json:"name"`
	Category       string           `json:"category"`
	Cooking_method string           `json:"cooking_method"`
	Composition    menu.Composition `json:"composition"`
	OrderIceType   string           `json:"order_ice_type"`
	Price          uint32           `json:"price"`
	BarsID         []uint32         `json:"bars_id,omitempty"`
}
