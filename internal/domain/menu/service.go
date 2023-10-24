package menu

import (
	"context"
	"restapi/pkg/logging"
)

type Service interface {
	NewMenu(context.Context, CreateMenuDTO) (string, error)
	DeleteMenu(context.Context, DeleteMenuDTO) error
	FindMenu(context.Context, FindMenuDTO) (Menu, error)
	FindUserMenus(context.Context, UserMenusDTO) (RespUserMenus, error)
	UpdateMenu(context.Context, UpdateMenuDTO) error
	UpdateMenuName(context.Context, UpdateMenuNameDTO) error
	AddDrink(context.Context, AddDrinkDTO) (string, error)
	DeleteDrink(context.Context, DeleteDrinkDTO) error
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

func (s *service) NewMenu(ctx context.Context, dto CreateMenuDTO) (string, error) {
	s.logger.Infof("creating menu %s", dto.Name)

	totalCost := s.GetTotalCost(dto.Drinks)

	menuID, err := s.repository.CreateMenu(ctx, dto, totalCost)

	if err != nil {
		return "", err
	}

	s.logger.Infof("menu is created, menu_id: %s", menuID)

	return menuID, nil
}

func (s *service) DeleteMenu(ctx context.Context, dto DeleteMenuDTO) error {
	s.logger.Infof("deleting menu %s", dto.ID)

	err := s.repository.DeleteMenu(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("menu is deleted, menu_id: %s", dto.ID)

	return nil
}

func (s *service) FindMenu(ctx context.Context, dto FindMenuDTO) (Menu, error) {
	s.logger.Infof("find user menu, menu_id: %s", dto.ID)

	mn, err := s.repository.FindMenu(ctx, dto)

	if err != nil {
		return Menu{}, err
	}

	s.logger.Infof("user menu is found")
	s.logger.Tracef("menu name: %s", mn.Name)

	return mn, nil
}

func (s *service) FindUserMenus(ctx context.Context, dto UserMenusDTO) (RespUserMenus, error) {
	s.logger.Infof("find user menus, user_id: %s", dto.UserID)

	menus, err := s.repository.FindUserMenus(ctx, dto)

	if err != nil {
		return RespUserMenus{}, err
	}

	s.logger.Infof("user menus are found")

	return menus, nil
}

func (s *service) UpdateMenu(ctx context.Context, dto UpdateMenuDTO) error {
	s.logger.Infof("update menu")

	totalCost := s.UpdateTotalCost(dto.Drinks)

	updatedID, err := s.repository.UpdateMenu(ctx, dto, totalCost)

	if err != nil {
		return err
	}

	s.logger.Infof("menu %s is updated", updatedID)

	return nil
}

func (s *service) UpdateMenuName(ctx context.Context, dto UpdateMenuNameDTO) error {
	s.logger.Infof("update menu name")

	err := s.repository.UpdateNameMenu(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("menu %s name is updated", dto.ID)

	return nil
}

func (s *service) AddDrink(ctx context.Context, dto AddDrinkDTO) (string, error) {
	s.logger.Infof("adding drink to menu %s", dto.MenuID)

	drinkID, err := s.repository.AddDrink(ctx, dto)

	if err != nil {
		return "", err
	}

	s.logger.Infof("drink added to menu")

	return drinkID, nil
}

func (s *service) DeleteDrink(ctx context.Context, dto DeleteDrinkDTO) error {
	s.logger.Infof("deleting drink from menu")

	err := s.repository.DeleteDrink(ctx, dto)
	if err != nil {
		return err
	}

	s.logger.Infof("drink deleted from menu")

	return nil
}

func (s *service) UpdateTotalCost(drinkGroups map[string][]Drink) uint32 {
	var totalCost uint32
	for _, drinks := range drinkGroups {
		for _, drink := range drinks {
			totalCost += drink.Price
		}
	}
	return totalCost
}

func (s *service) GetTotalCost(drinkGroups map[string][]NewDrinkDTO) uint32 {
	var totalCost uint32
	for _, drinks := range drinkGroups {
		for _, drink := range drinks {
			totalCost += drink.Price
		}
	}
	return totalCost
}
