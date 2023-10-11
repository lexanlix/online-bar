package bar_api

import (
	"context"
	"encoding/json"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/bar"
	"restapi/internal/domain/user"
	"strconv"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	createBarURL = "/api/bar/create"
	closeBarURL  = "/api/bar/close"
	updateBarURL = "/api/bar/update"

	getBarOrdersURL = "/api/bar/orders"
)

type handler struct {
	userService user.Service
	service     bar.Service
	logger      *logging.Logger
}

func NewHandler(logger *logging.Logger, service bar.Service, userService user.Service) adapters.Handler {
	return &handler{
		service:     service,
		userService: userService,
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createBarURL, apperror.Middleware(h.Verify(h.CreateBar)))
	router.HandlerFunc(http.MethodDelete, closeBarURL, apperror.Middleware(h.Verify(h.CloseBar)))
	router.HandlerFunc(http.MethodGet, getBarOrdersURL, apperror.Middleware(h.Verify(h.GetBarOrders)))
	router.HandlerFunc(http.MethodPost, updateBarURL, apperror.Middleware(h.Verify(h.UpdateBar)))
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

func (h *handler) GetBarOrders(w http.ResponseWriter, r *http.Request) error {
	var dto bar.GetBarOrdersDTO

	id, err := strconv.Atoi(r.URL.Query().Get("bar_id"))
	if err != nil {
		return apperror.NewAppError(err, "wrong data", err.Error(), "US-000009")
	}

	dto.ID = uint32(id)
	dto.EventID = r.URL.Query().Get("event_id")

	if dto.EventID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param event_id is empty", "US-000015")
	}

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
