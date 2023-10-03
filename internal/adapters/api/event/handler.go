package event_api

import (
	"context"
	"encoding/json"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/event"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	createEventURL     = "/api/event/create"
	deleteEventURL     = "/api/event/delete"
	getEventsByHostURL = "/api/user/events"
	getEventByIDurl    = "/api/user/event"
	updateEventURL     = "/api/event/update"
)

type handler struct {
	service event.Service
	logger  *logging.Logger
}

func NewHandler(logger *logging.Logger, service event.Service) adapters.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createEventURL, apperror.Middleware(h.CreateEvent))
	router.HandlerFunc(http.MethodDelete, deleteEventURL, apperror.Middleware(h.DeleteEvent))
	router.HandlerFunc(http.MethodPut, getEventsByHostURL, apperror.Middleware(h.GetAllByHostID))
	router.HandlerFunc(http.MethodPut, getEventByIDurl, apperror.Middleware(h.GetByID))
	router.HandlerFunc(http.MethodPost, updateEventURL, apperror.Middleware(h.UpdateEvent))

	// Обработчики, доступные пользователям, вошедшим в аккаунт (которые имеют AccessToken)
	//router.HandlerFunc(http.MethodDelete, deleteUserURL, apperror.Middleware(h.Verify(h.DeleteUser)))
}

func (h *handler) CreateEvent(w http.ResponseWriter, r *http.Request) error {
	var dto event.CreateEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	eventID, err2 := h.service.NewEvent(context.TODO(), dto)
	if err != nil {
		return err2
	}

	resp := event.RespCreateEvent{
		ID: eventID,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) DeleteEvent(w http.ResponseWriter, r *http.Request) error {
	var dto event.DeleteEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.DeleteEvent(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("event is deleted"))

	return nil
}

func (h *handler) GetAllByHostID(w http.ResponseWriter, r *http.Request) error {
	var dto event.FindAllEventsDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	allUserEvents, err := h.service.FindAllUserEvents(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	allBytes, err := json.Marshal(allUserEvents)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	var dto event.FindEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	userEvent, err := h.service.FindEvent(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	eventBytes, err := json.Marshal(userEvent)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(eventBytes)

	return nil
}

func (h *handler) UpdateEvent(w http.ResponseWriter, r *http.Request) error {
	var dto event.UpdateEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.UpdateEvent(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("event is updated"))

	return nil
}
