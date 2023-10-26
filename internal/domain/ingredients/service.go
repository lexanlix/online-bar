package ingredients

import (
	"context"
	"fmt"
	"restapi/pkg/logging"
)

type Service interface {
	NewIngredients(context.Context, AddIngredientsDTO) error
	AddIngredient(context.Context, AddIngredientDTO) error
	DeleteIngredient(context.Context, DeleteIngredientDTO) error
	DeleteEventIngredients(context.Context, DeleteEventIngrDTO) error
	FindIngredient(context.Context, FindIngredientDTO) (Ingredient, error)
	FindEventIngredients(context.Context, FindEventIngredientsDTO) (RespEventIngredients, error)
	UpdateIngredient(context.Context, UpdateIngredientDTO) error
	Validate(AddIngredientsDTO) error
}

type service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) NewIngredients(ctx context.Context, dto AddIngredientsDTO) error {
	s.logger.Infof("creating list of event ingredients")

	IDs, err := s.repository.AddIngredients(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("list of ingredients created:")
	for i, ingrID := range IDs {
		s.logger.Tracef("\n%d id: %s", i, ingrID)
	}

	return nil
}

func (s *service) AddIngredient(ctx context.Context, dto AddIngredientDTO) error {
	s.logger.Infof("adding ingredient %s", dto.Name)

	ingrID, err := s.repository.AddIngredient(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("ingredient %s added to list", ingrID)

	return nil
}

func (s *service) DeleteIngredient(ctx context.Context, dto DeleteIngredientDTO) error {
	s.logger.Infof("deleting ingredient")

	err := s.repository.DeleteIngredient(ctx, dto)
	if err != nil {
		return err
	}

	s.logger.Infof("ingredient deleted")

	return nil
}

func (s *service) DeleteEventIngredients(ctx context.Context, dto DeleteEventIngrDTO) error {
	s.logger.Infof("deleting event %s ingredients", dto.EventID)

	err := s.repository.DeleteEventIngredients(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("event %s ingredients deleted", dto.EventID)

	return nil
}

func (s *service) FindIngredient(ctx context.Context, dto FindIngredientDTO) (Ingredient, error) {
	s.logger.Infof("find ingredient %s", dto.ID)

	ingr, err := s.repository.FindIngredient(ctx, dto)

	if err != nil {
		return Ingredient{}, err
	}

	s.logger.Infof("ingredient is found")

	return ingr, nil
}

func (s *service) FindEventIngredients(ctx context.Context, dto FindEventIngredientsDTO) (RespEventIngredients, error) {
	s.logger.Infof("find event ingredients, event_id: %s", dto.EventID)

	var resp RespEventIngredients
	var err error

	resp.Ingredients, err = s.repository.FindEventIngredients(ctx, dto)

	if err != nil {
		return RespEventIngredients{}, err
	}

	s.logger.Infof("event ingredients are found")

	return resp, nil
}

func (s *service) UpdateIngredient(ctx context.Context, dto UpdateIngredientDTO) error {
	s.logger.Infof("update ingredient")

	updatedID, err := s.repository.UpdateIngredient(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("ingredient %s is updated", updatedID)

	return nil
}

// TODO
func (s *service) Validate(dto AddIngredientsDTO) error {

	// Добавить проверку на то, что список ингредиентов для данного ивента уже создан

	if dto.UserID == "" {
		return fmt.Errorf("user id field is empty")
	}

	if dto.EventID == "" {
		return fmt.Errorf("event id field is empty")
	}

	if len(dto.Ingredients) == 0 {
		return fmt.Errorf("ingredients list is empty")
	}

	return nil
}
