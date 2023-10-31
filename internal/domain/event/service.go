package event

import (
	"context"
	"fmt"
	"restapi/internal/domain/menu"
	"restapi/pkg/logging"
	"time"
)

const (
	statusCreated   = "Created"
	statusActive    = "Active"
	statusCompleted = "Completed"
)

type Service interface {
	NewEvent(context.Context, CreateEventDTO) (Event, error)
	SetActive(timer *time.Timer, id string)
	CompleteEvent(context.Context, CompleteEventDTO) error
	FindAllUserEvents(context.Context, FindAllEventsDTO) ([]Event, error)
	FindEvent(context.Context, FindEventDTO) (Event, error)
	UpdateEvent(context.Context, UpdateEventDTO) error
}

type service struct {
	repository Repository
	menuRepos  menu.Repository
	logger     *logging.Logger
}

func NewService(repository Repository, menuRepos menu.Repository, logger *logging.Logger) Service {
	return &service{
		repository: repository,
		menuRepos:  menuRepos,
		logger:     logger,
	}
}

func (s *service) NewEvent(ctx context.Context, dto CreateEventDTO) (Event, error) {
	s.logger.Infof("creating event %s", dto.Name)

	// Так как time.Now() дает время с часовым поясом, а в DateTime без, то вычитаем 3 часа
	timeUntil := (time.Until(dto.DateTime) - (3 * time.Hour))
	timer := time.NewTimer(timeUntil)

	// Получить напитки из меню и найти наименования всех ингредиентов
	menuDTO := menu.FindMenuDTO{
		ID: dto.MenuID,
	}

	menu, err := s.menuRepos.FindMenu(ctx, menuDTO)
	if err != nil {
		return Event{}, fmt.Errorf("finding menu error: %v", err)
	}

	shopList := s.GetShoppingList(menu)

	eventID, err := s.repository.CreateEvent(ctx, dto, shopList)

	if err != nil {
		return Event{}, err
	}

	evnt := Event{
		ID:                 eventID,
		HostID:             dto.HostID,
		Name:               dto.Name,
		Description:        dto.Description,
		ParticipantsNumber: dto.ParticipantsNumber,
		DateTime:           dto.DateTime,
		Status:             statusCreated,
		MenuID:             dto.MenuID,
		ShoppingList:       shopList,
	}

	go s.SetActive(timer, eventID)

	s.logger.Infof("event is created, event_id: %s", eventID)

	return evnt, nil
}

// Мероприятие становится активным
func (s *service) SetActive(timer *time.Timer, id string) {
	<-timer.C

	status, err := s.repository.SetActive(context.TODO(), id)
	if err != nil {
		panic(err)
	}

	s.logger.Infof("event %s now is %s", id, status)
}

func (s *service) CompleteEvent(ctx context.Context, dto CompleteEventDTO) error {
	s.logger.Infof("completing event %s", dto.ID)

	status, err := s.repository.DeleteEvent(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("event is %s, event_id: %s", status, dto.ID)

	return nil
}

func (s *service) FindAllUserEvents(ctx context.Context, dto FindAllEventsDTO) ([]Event, error) {
	s.logger.Infof("find all user events, user_id: %s", dto.HostID)

	events, err := s.repository.FindAllUserEvents(ctx, dto)

	if err != nil {
		return nil, err
	}

	s.logger.Infof("all user events is found")
	for n, evnt := range events {
		s.logger.Tracef("\n%d event name: %s", n, evnt.Name)
	}

	return events, nil
}

func (s *service) FindEvent(ctx context.Context, dto FindEventDTO) (Event, error) {
	s.logger.Infof("find one user event, user_id: %s, event_id: %s", dto.HostID, dto.ID)

	evnt, err := s.repository.FindUserEvent(ctx, dto)

	if err != nil {
		return Event{}, err
	}

	s.logger.Infof("user event is found")
	s.logger.Tracef("event name: %s", evnt.Name)

	return evnt, nil
}

func (s *service) UpdateEvent(ctx context.Context, dto UpdateEventDTO) error {
	s.logger.Infof("update event")

	updatedID, err := s.repository.UpdateEvent(ctx, dto)

	if err != nil {
		return err
	}

	s.logger.Infof("event %s is updated", updatedID)

	return nil
}

// Может быть как то оптимизировать?
func (s *service) GetShoppingList(menu menu.Menu) []string {
	Hash := make(map[string]bool, 0)

	for _, drinks := range menu.Drinks {
		for _, drink := range drinks {
			for _, liq := range drink.Composition.Liquids {
				Hash[liq.Name] = true
			}

			for _, solB := range drink.Composition.SolidsBulk {
				Hash[solB.Name] = true
			}

			for _, solU := range drink.Composition.SolidsUnit {
				Hash[solU.Name] = true
			}

			switch drink.OrderIceType {
			case "block_ice":
				Hash["Лед блоками (block)"] = true
			case "cubed_ice":
				Hash["Лед кубиками (cubed)"] = true
			case "cracked_ice":
				Hash["Ломаный лед (cracked)"] = true
			case "nugget_ice":
				Hash["Пальчиковый лед (nugget)"] = true
			case "crushed_ice":
				Hash["Дробленый лед (crushed)"] = true
			}
		}
	}

	list := make([]string, 0, len(Hash))
	for k := range Hash {
		list = append(list, k)
	}

	return list
}
