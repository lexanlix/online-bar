package event

import (
	"context"
	"restapi/pkg/logging"
	"time"
)

type Service interface {
	NewEvent(context.Context, CreateEventDTO) (string, error)
	SetActive(timer *time.Timer, id string)
	CompleteEvent(context.Context, CompleteEventDTO) error
	FindAllUserEvents(context.Context, FindAllEventsDTO) ([]Event, error)
	FindEvent(context.Context, FindEventDTO) (Event, error)
	UpdateEvent(context.Context, UpdateEventDTO) error

	//SetMenu(context.Context, menu.CreateMenuDTO) (menu.Menu, error)
	//UpdateMenu(context.Context) error
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

func (s *service) NewEvent(ctx context.Context, dto CreateEventDTO) (string, error) {
	s.logger.Infof("creating event %s", dto.Name)

	// Так как time.Now() дает время с часовым поясом, а в DateTime без, то вычитаем 3 часа
	timeUntil := (time.Until(dto.DateTime) - (3 * time.Hour))
	timer := time.NewTimer(timeUntil)

	eventID, err := s.repository.CreateEvent(ctx, dto)

	if err != nil {
		return "", err
	}

	go s.SetActive(timer, eventID)

	s.logger.Infof("event is created, event_id: %s", eventID)

	return eventID, nil
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

	evnt, err := s.repository.FindOneUserEvent(ctx, dto)

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

// // Править
// func (s *service) SetMenu(ctx context.Context, dto menu.CreateMenuDTO) (menu.Menu, error) {
// 	newMenu := menu.NewMenu(dto.ID, dto.Name, dto.Drinks)
// 	newMenu.UpdateTotalCost()

// 	return newMenu, nil
// }

// // Править
// func (s *service) UpdateMenu(context.Context) error {
// 	panic("TODO this")
// }
