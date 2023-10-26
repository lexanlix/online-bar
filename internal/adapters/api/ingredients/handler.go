package ingredients_api

import (
	"context"
	"encoding/json"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/ingredients"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	newIngrListURL    = "/api/event/ingr/new_list"
	addIngrURL        = "/api/event/ingr/add"
	getIngrListURL    = "/api/event/ingr_list"
	getIngrURL        = "/api/event/ingr"
	deleteIngrListURL = "/api/event/ingr/delete_list"
	deleteIngrURL     = "/api/event/ingr/delete"
	updateIngrURL     = "/api/event/ingr/update"
)

type handler struct {
	service ingredients.Service
	logger  *logging.Logger
}

func NewHandler(logger *logging.Logger, service ingredients.Service) adapters.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, newIngrListURL, apperror.Middleware(h.NewIngrList))
	router.HandlerFunc(http.MethodPost, addIngrURL, apperror.Middleware(h.AddIngr))
	router.HandlerFunc(http.MethodGet, getIngrListURL, apperror.Middleware(h.GetIngrList))
	router.HandlerFunc(http.MethodGet, getIngrURL, apperror.Middleware(h.GetIngr))
	router.HandlerFunc(http.MethodDelete, deleteIngrListURL, apperror.Middleware(h.DeleteIngrList))
	router.HandlerFunc(http.MethodDelete, deleteIngrURL, apperror.Middleware(h.DeleteIngr))
	router.HandlerFunc(http.MethodPut, updateIngrURL, apperror.Middleware(h.UpdateIngr))
}

func (h *handler) NewIngrList(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.AddIngredientsDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.Validate(dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong ingredient list add data", err.Error(), "US-000009")
	}

	err = h.service.NewIngredients(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong ingredient list add data", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) AddIngr(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.AddIngredientDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.AddIngredient(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong ingredient add data", err.Error(), "US-000009")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) GetIngrList(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.FindEventIngredientsDTO
	dto.EventID = r.URL.Query().Get("event_id")

	if dto.EventID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param event_id is empty", "US-000015")
	}

	resp, err := h.service.FindEventIngredients(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong event data", err.Error(), "US-000009")
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(respBytes)

	return nil
}

func (h *handler) GetIngr(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.FindIngredientDTO
	dto.ID = r.URL.Query().Get("id")

	if dto.ID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param id is empty", "US-000015")
	}

	resp, err := h.service.FindIngredient(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong ingredient data", err.Error(), "US-000009")
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(respBytes)

	return nil
}

func (h *handler) DeleteIngrList(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.DeleteEventIngrDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.DeleteEventIngredients(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong event id", err.Error(), "US-000009")
	}

	w.WriteHeader(200)
	w.Write([]byte("ingredients list is deleted"))

	return nil
}

func (h *handler) DeleteIngr(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.DeleteIngredientDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.DeleteIngredient(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(200)
	w.Write([]byte("ingredient is deleted"))

	return nil
}

func (h *handler) UpdateIngr(w http.ResponseWriter, r *http.Request) error {
	var dto ingredients.UpdateIngredientDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.UpdateIngredient(context.Background(), dto)
	if err != nil {
		return apperror.NewAppError(err, "wrong update ingredient data", err.Error(), "US-000009")
	}

	w.WriteHeader(204)

	return nil
}
