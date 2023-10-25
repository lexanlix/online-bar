package ingredients

import "context"

type Repository interface {
	AddIngredient(context.Context, AddIngredientDTO) (string, error)
	AddIngredients(context.Context, AddIngredientsDTO) ([]string, error)
	DeleteIngredient(context.Context, DeleteIngredientDTO) error
	DeleteEventIngredients(context.Context, DeleteEventIngrDTO) error
	FindIngredient(context.Context, FindIngredientDTO) (Ingredient, error)
	FindEventIngredients(context.Context, FindEventIngredientsDTO) ([]Ingredient, error)
	UpdateIngredient(context.Context, UpdateIngredientDTO) (string, error)
}
