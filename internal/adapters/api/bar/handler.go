package bar_api

import (
	"context"
	"encoding/json"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/bar"
	"strconv"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	createBarURL      = "/api/bar/create"
	closeBarURL       = "/api/bar/close"
	getEventOrdersURL = "/api/event/orders"
	getBarOrdersURL   = "/api/event/bar/orders"
	updateBarURL      = "/api/bar/update"
)

type handler struct {
	service bar.Service
	logger  *logging.Logger
}

func NewHandler(logger *logging.Logger, service bar.Service) adapters.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	//router.HandlerFunc(http.MethodPost, createBarURL, apperror.Middleware(h.CreateBar))
	router.HandlerFunc(http.MethodDelete, closeBarURL, apperror.Middleware(h.CloseBar))
	router.HandlerFunc(http.MethodGet, getEventOrdersURL, apperror.Middleware(h.GetEventOrders))
	router.HandlerFunc(http.MethodGet, getBarOrdersURL, apperror.Middleware(h.GetBarOrders))
	router.HandlerFunc(http.MethodPost, updateBarURL, apperror.Middleware(h.UpdateBar))

	// Обработчики, доступные пользователям, вошедшим в аккаунт (которые имеют AccessToken)
}

func (h *handler) CreateBar(w http.ResponseWriter, r *http.Request) error {
	var dto bar.CreateBarDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	barID, menu, err2 := h.service.OpenBar(context.TODO(), dto)
	if err2 != nil {
		return err2
	}

	resp := bar.RespCreateBar{
		ID:   barID,
		Menu: menu,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) CloseBar(w http.ResponseWriter, r *http.Request) error {
	var dto bar.CloseBarDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.CloseBar(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("bar is closed"))

	return nil
}

func (h *handler) GetEventOrders(w http.ResponseWriter, r *http.Request) error {
	var dto bar.GetOrdersDTO

	dto.EventID = r.Header.Get("event_id")

	allEventOrders, err := h.service.GetOrders(context.TODO(), dto)
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

func (h *handler) GetBarOrders(w http.ResponseWriter, r *http.Request) error {
	var dto bar.GetBarOrdersDTO

	id, err := strconv.Atoi(r.Header.Get("bar_id"))
	if err != nil {
		return apperror.NewAppError(err, "wrong data", err.Error(), "US-000009")
	}

	dto.ID = uint32(id)
	dto.EventID = r.Header.Get("event_id")

	barOrders, err := h.service.GetBarOrders(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	ordersBytes, err := json.Marshal(barOrders)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ordersBytes)

	return nil
}

func (h *handler) UpdateBar(w http.ResponseWriter, r *http.Request) error {
	var dto bar.UpdateBarDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.UpdateInfo(context.TODO(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("bar info is updated"))

	return nil
}
