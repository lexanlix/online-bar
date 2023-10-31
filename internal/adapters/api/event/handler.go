package event_api

import (
	"context"
	"encoding/json"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/bar"
	"restapi/internal/domain/event"
	"restapi/internal/domain/user"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	createEventURL = "/api/event/create"

	getEventsByHostURL = "/api/user/events"
	getEventByIDurl    = "/api/user/event"
	getEventOrdersURL  = "/api/event/orders"

	completeEventURL = "/api/event/complete"
	updateEventURL   = "/api/event/update"
)

type handler struct {
	service     event.Service
	userService user.Service
	barService  bar.Service
	logger      *logging.Logger
}

func NewHandler(logger *logging.Logger, service event.Service, userService user.Service,
	barService bar.Service) adapters.Handler {
	return &handler{
		service:     service,
		userService: userService,
		barService:  barService,
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createEventURL, apperror.Middleware(h.CreateEvent))
	router.HandlerFunc(http.MethodPatch, completeEventURL, apperror.Middleware(h.CompleteEvent))
	router.HandlerFunc(http.MethodGet, getEventsByHostURL, apperror.Middleware(h.GetAllByHostID))
	router.HandlerFunc(http.MethodGet, getEventByIDurl, apperror.Middleware(h.GetByID))
	router.HandlerFunc(http.MethodGet, getEventOrdersURL, apperror.Middleware(h.GetEventOrders))
	router.HandlerFunc(http.MethodPut, updateEventURL, apperror.Middleware(h.UpdateEvent))
}

func (h *handler) CreateEvent(w http.ResponseWriter, r *http.Request) error {
	var dto event.CreateEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	event, err := h.service.NewEvent(context.TODO(), dto)
	if err != nil {
		return err
	}

	respBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) CompleteEvent(w http.ResponseWriter, r *http.Request) error {
	var dto event.CompleteEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.CompleteEvent(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("event is completed"))

	return nil
}

func (h *handler) GetAllByHostID(w http.ResponseWriter, r *http.Request) error {
	var dto event.FindAllEventsDTO

	dto.UserID = r.URL.Query().Get("user_id")

	if dto.UserID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param user_id is empty", "US-000015")
	}

	resp, err := h.service.FindAllUserEvents(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	var dto event.FindEventDTO

	dto.ID = r.URL.Query().Get("event_id")
	dto.UserID = r.URL.Query().Get("user_id")

	if dto.UserID == "" || dto.ID == "" {
		if dto.UserID == "" {
			return apperror.NewAppError(nil, "query param is empty", "param user_id is empty", "US-000015")
		} else {
			return apperror.NewAppError(nil, "query param is empty", "param event_id is empty", "US-000015")
		}
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

// TODO
func (h *handler) GetEventOrders(w http.ResponseWriter, r *http.Request) error {
	var dto bar.GetOrdersDTO

	dto.EventID = r.URL.Query().Get("event_id")

	if dto.EventID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param event_id is empty", "US-000015")
	}

	allEventOrders, err := h.barService.GetOrders(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	allBytes, err := json.Marshal(allEventOrders)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

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

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) Verify(protectedHandler apperror.AppHandler) apperror.AppHandler {

	return func(w http.ResponseWriter, r *http.Request) error {
		cookie, err := r.Cookie("AccessToken")
		if err != nil {
			h.logger.Errorf("cookie error: %v", err)
			return apperror.ErrUnauthorized
		}

		if cookie.Value == "" {
			h.logger.Errorf("access token is empty")
			return apperror.ErrUnauthorized
		}

		err = h.userService.Verify(context.TODO(), cookie.Value)
		if err != nil {
			h.logger.Errorf("access token is wrong: %v", err)
			return apperror.ErrUnauthorized
		}

		return protectedHandler(w, r)
	}
}
