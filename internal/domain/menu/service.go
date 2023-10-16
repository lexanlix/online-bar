package menu

import (
	"context"
	"restapi/pkg/logging"
)

type Service interface {
	NewMenu(context.Context, CreateMenuDTO) (string, error)
	DeleteMenu(context.Context, DeleteMenuDTO) error
	FindMenu(context.Context, FindMenuDTO) (Menu, error)
	UpdateMenu(context.Context, UpdateMenuDTO) error
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

	menuID, err := s.repository.CreateMenu(ctx, dto)

	if err != nil {
		return "", err
	}

	s.logger.Infof("menu is created, menu_id: %s", menuID)

	return menuID, nil
}

func (s *service) DeleteMenu(ctx context.Context, dto DeleteMenuDTO) error {
	s.logger.Infof("deleting menu %s", dto.ID)

	status, err := s.repository.DeleteMenu(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("menu is %s, menu_id: %s", status, dto.ID)

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

func (s *service) UpdateMenu(ctx context.Context, dto UpdateMenuDTO) error {
	s.logger.Infof("update menu")

	updatedID, err := s.repository.UpdateMenu(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("menu %s is updated", updatedID)

	return nil
}
