package drinks_list

import (
	"context"
	"restapi/internal/domain/menu"
	"restapi/pkg/logging"
)

type Service interface {
	AddUserDrink(context.Context, AddUserDrinkDTO) (menu.Drink, error)
	DeleteUserDrink(context.Context, DeleteUserDrinkDTO) error
	FindUserDrink(context.Context, FindUserDrinkDTO) (menu.Drink, error)
	FindUserDrinks(context.Context, FindUserDrinksDTO) (RespFindUDrinks, error)
	UpdateUserDrink(context.Context, UpdateUserDrinkDTO) error
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

func (s *service) AddUserDrink(ctx context.Context, dto AddUserDrinkDTO) (menu.Drink, error) {
	s.logger.Infof("adding drink to drink list")

	drinkID, err := s.repository.AddUserDrink(ctx, dto)

	if err != nil {
		return menu.Drink{}, err
	}

	dr := menu.Drink{
		ID:             drinkID,
		Name:           dto.Name,
		Category:       dto.Category,
		Cooking_method: dto.Cooking_method,
		Composition:    dto.Composition,
		OrderIceType:   dto.OrderIceType,
		Price:          dto.Price,
		BarsID:         dto.BarsID,
	}

	s.logger.Infof("drink added to drink list")

	return dr, nil
}

func (s *service) DeleteUserDrink(ctx context.Context, dto DeleteUserDrinkDTO) error {
	s.logger.Infof("deleting drink from drinks list")

	err := s.repository.DeleteUserDrink(ctx, dto)
	if err != nil {
		return err
	}

	s.logger.Infof("drink deleted from drinks list")

	return nil
}

func (s *service) FindUserDrink(ctx context.Context, dto FindUserDrinkDTO) (menu.Drink, error) {
	s.logger.Infof("find user drink")

	UserDrink, err := s.repository.FindUserDrink(ctx, dto)

	if err != nil {
		return menu.Drink{}, err
	}

	s.logger.Infof("user drink is found")
	s.logger.Tracef("drink name: %s", UserDrink.Name)

	return UserDrink, nil
}

func (s *service) FindUserDrinks(ctx context.Context, dto FindUserDrinksDTO) (RespFindUDrinks, error) {
	s.logger.Infof("find user drinks, user_id: %s", dto.UserID)

	menus, err := s.repository.FindUserDrinks(ctx, dto)

	if err != nil {
		return RespFindUDrinks{}, err
	}

	s.logger.Infof("user drinks are found")

	return menus, nil
}

func (s *service) UpdateUserDrink(ctx context.Context, dto UpdateUserDrinkDTO) error {
	s.logger.Infof("update user drink")

	updatedID, err := s.repository.UpdateUserDrink(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("user drink %s is updated", updatedID)

	return nil
}
