package drinks_list_api

import (
	"context"
	"encoding/json"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/drinks_list"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	addUserDrinkURL    = "/api/user/drinks/add"
	deleteUserDrinkURL = "/api/user/drinks/delete"
	getUserDrinkURL    = "/api/user/drink"
	getUserDrinksURL   = "/api/user/drinks"
	updateUserDrinkURL = "/api/user/drinks/update"
)

type handler struct {
	service drinks_list.Service
	logger  *logging.Logger
}

func NewHandler(logger *logging.Logger, service drinks_list.Service) adapters.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, addUserDrinkURL, apperror.Middleware(h.AddUserDrink))
	router.HandlerFunc(http.MethodDelete, deleteUserDrinkURL, apperror.Middleware(h.DeleteUserDrink))
	router.HandlerFunc(http.MethodGet, getUserDrinkURL, apperror.Middleware(h.GetUserDrink))
	router.HandlerFunc(http.MethodGet, getUserDrinksURL, apperror.Middleware(h.GetUserDrinks))
	router.HandlerFunc(http.MethodPut, updateUserDrinkURL, apperror.Middleware(h.UpdateUserDrink))
}

func (h *handler) AddUserDrink(w http.ResponseWriter, r *http.Request) error {
	var dto drinks_list.AddUserDrinkDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	drink, err := h.service.AddUserDrink(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong user drink add data", err.Error(), "US-000009")
	}

	respBytes, err := json.Marshal(drink)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) DeleteUserDrink(w http.ResponseWriter, r *http.Request) error {
	var dto drinks_list.DeleteUserDrinkDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.DeleteUserDrink(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(200)
	w.Write([]byte("user drink is deleted"))

	return nil
}

func (h *handler) GetUserDrink(w http.ResponseWriter, r *http.Request) error {
	var dto drinks_list.FindUserDrinkDTO
	dto.ID = r.URL.Query().Get("id")

	if dto.ID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param id is empty", "US-000015")
	}

	resp, err := h.service.FindUserDrink(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong user drink data", err.Error(), "US-000009")
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(respBytes)

	return nil
}

func (h *handler) GetUserDrinks(w http.ResponseWriter, r *http.Request) error {
	var dto drinks_list.FindUserDrinksDTO
	dto.UserID = r.URL.Query().Get("user_id")

	if dto.UserID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param user_id is empty", "US-000015")
	}

	resp, err := h.service.FindUserDrinks(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong user data", err.Error(), "US-000009")
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(respBytes)

	return nil
}

func (h *handler) UpdateUserDrink(w http.ResponseWriter, r *http.Request) error {
	var dto drinks_list.UpdateUserDrinkDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.UpdateUserDrink(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong update user drink data", err.Error(), "US-000009")
	}

	w.WriteHeader(204)

	return nil
}
