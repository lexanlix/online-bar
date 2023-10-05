package bar

import (
	"context"
	"restapi/internal/domain/menu"
	"restapi/pkg/logging"
)

type Service interface {
	OpenBar(context.Context, CreateBarDTO) (uint32, menu.Menu, error)
	CloseBar(context.Context, CloseBarDTO) error
	UpdateInfo(context.Context, UpdateBarDTO) error
	GetOrders(context.Context, GetOrdersDTO) ([]string, error)
	GetBarOrders(context.Context, GetBarOrdersDTO) ([]string, error)
	SetMenu(context.Context, menu.CreateMenuDTO) (menu.Menu, error)
	UpdateMenu(context.Context) error
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

func (s *service) OpenBar(ctx context.Context, dto CreateBarDTO) (uint32, menu.Menu, error) {
	s.logger.Infof("creating bar")

	barID, err := s.repository.CreateBar(ctx, dto)

	if err != nil {
		return 0, menu.Menu{}, err
	}

	mn, err := s.SetMenu(ctx, dto.CreateMenuDTO)
	if err != nil {
		return 0, menu.Menu{}, err
	}

	s.logger.Infof("bar is created, bar_id: %d", barID)

	return barID, mn, nil
}

func (s *service) CloseBar(ctx context.Context, dto CloseBarDTO) error {
	s.logger.Infof("closing bar %d", dto.ID)

	status, err := s.repository.CloseBar(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("bar is %s, bar_id: %d", status, dto.ID)

	return nil
}

func (s *service) GetOrders(ctx context.Context, dto GetOrdersDTO) ([]string, error) {
	s.logger.Infof("find orders from all event bars, event_id: %s", dto.EventID)

	ordersID, err := s.repository.GetOrders(ctx, dto)

	if err != nil {
		return nil, err
	}

	s.logger.Infof("all event orders is found")
	for n, id := range ordersID {
		s.logger.Tracef("\n%d event name: %s", n, id)
	}

	return ordersID, nil
}

func (s *service) GetBarOrders(ctx context.Context, dto GetBarOrdersDTO) ([]string, error) {
	s.logger.Infof("find one bar orders, event_id: %s, bar_id: %d", dto.EventID, dto.ID)

	ordersID, err := s.repository.GetBarOrders(ctx, dto)

	if err != nil {
		return nil, err
	}

	s.logger.Infof("event orders is found")
	s.logger.Tracef("orders id: %s", ordersID)

	return ordersID, nil
}

func (s *service) UpdateInfo(ctx context.Context, dto UpdateBarDTO) error {
	s.logger.Infof("update bar")

	updatedID, err := s.repository.UpdateInfo(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("bar %s is updated", updatedID)

	return nil
}

func (s *service) SetMenu(ctx context.Context, dto menu.CreateMenuDTO) (menu.Menu, error) {
	newMenu := menu.NewMenu(dto.Drinks)
	newMenu.UpdateTotalCost()

	return newMenu, nil
}

func (s *service) UpdateMenu(context.Context) error {
	panic("TODO this")
}
