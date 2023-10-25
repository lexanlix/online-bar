package ingredients

const (
	// ice_types
	BlockIce   = "block_ice"
	CubedIce   = "cubed_ice"
	CrackedIce = "cracked_ice"
	NuggetIce  = "nugget_ice"
	CrushedIce = "crushed_ice"
	NoIce      = "no_ice"

	// ingredient_types
	IceType       = "ice"
	LiquidType    = "liquids"
	SolidBulkType = "solids_bulk"
	SolidUnitType = "solids_unit"
)

// type - тип ингридиента
// У типа IceType поле Name содержит ice_type
// Ингридиентов с IceType типом может быть несколько с уникальным name
// У ингридиентов типа SolidUnitType поле unit - пустое?, а в поле volume(объем) указывается количество в штуках
type Ingredient struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Unit    string `json:"unit,omitempty"`
	Volume  uint32 `json:"volume"`
	Cost    uint32 `json:"cost"`
}
