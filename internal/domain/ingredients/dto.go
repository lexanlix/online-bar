package ingredients

type AddIngredientsDTO struct {
	UserID      string              `json:"user_id"`
	EventID     string              `json:"event_id"`
	Ingredients []IngredientDataDTO `json:"ingredients"`
}

type IngredientDataDTO struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Unit   string `json:"unit,omitempty"`
	Volume uint32 `json:"volume"`
	Cost   uint32 `json:"cost"`
}

type AddIngredientDTO struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Unit    string `json:"unit,omitempty"`
	Volume  uint32 `json:"volume"`
	Cost    uint32 `json:"cost"`
}

type DeleteIngredientDTO struct {
	ID string `json:"id"`
}

type DeleteEventIngrDTO struct {
	EventID string `json:"event_id"`
}

type FindIngredientDTO struct {
	ID string `json:"id"`
}

type FindEventIngredientsDTO struct {
	EventID string `json:"event_id"`
}

type UpdateIngredientDTO struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Unit   string `json:"unit,omitempty"`
	Volume uint32 `json:"volume"`
	Cost   uint32 `json:"cost"`
}

type RespEventIngredients struct {
	Ingredients []Ingredient `json:"ingredients"`
}
