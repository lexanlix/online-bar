package user_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/bar"
	"restapi/internal/domain/event"
	"restapi/internal/domain/menu"
	"restapi/internal/domain/user"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	getUserURL     = "/api/user"
	signUpURL      = "/api/register"
	signInURL      = "/api/login"
	refreshURL     = "/api/auth/refresh"
	updateUserURL  = "/api/user/update"
	pUpdateUserURL = "/api/user/update/part"
	deleteUserURL  = "/api/user/delete"
	createEventURL = "/api/event/create"
	createBarURL   = "/api/bar/create"

	createMenuURL = "/api/user/menu/new"
	addDrinkURL   = "/api/user/menu/add"
)

type handler struct {
	eventService event.Service
	barService   bar.Service
	service      user.Service
	logger       *logging.Logger
}

func NewHandler(logger *logging.Logger, service user.Service, eventService event.Service,
	barService bar.Service) adapters.Handler {
	return &handler{
		service:      service,
		logger:       logger,
		eventService: eventService,
		barService:   barService,
	}
}

func (h *handler) Register(router *httprouter.Router) {

	router.HandlerFunc(http.MethodPost, signUpURL, apperror.Middleware(h.SignUp))
	router.HandlerFunc(http.MethodPost, signInURL, apperror.Middleware(h.SignIn))
	router.HandlerFunc(http.MethodGet, refreshURL, apperror.Middleware(h.UserRefresh))

	// for testing menu functions
	router.HandlerFunc(http.MethodPost, createMenuURL, apperror.Middleware(h.NewMenu))
	router.HandlerFunc(http.MethodPost, addDrinkURL, apperror.Middleware(h.AddDrink))

	// Обработчики, доступные пользователям, вошедшим в аккаунт (которые имеют AccessToken)
	router.HandlerFunc(http.MethodPatch, pUpdateUserURL, apperror.Middleware(h.Verify(h.PartiallyUpdateUser)))
	router.HandlerFunc(http.MethodGet, getUserURL, apperror.Middleware(h.Verify(h.GetUserByUUID)))
	router.HandlerFunc(http.MethodPut, updateUserURL, apperror.Middleware(h.Verify(h.UpdateUser)))
	router.HandlerFunc(http.MethodDelete, deleteUserURL, apperror.Middleware(h.Verify(h.DeleteUser)))
	router.HandlerFunc(http.MethodPost, createEventURL, apperror.Middleware(h.Verify(h.CreateEvent)))
	router.HandlerFunc(http.MethodPost, createBarURL, apperror.Middleware(h.Verify(h.CreateBar)))
}

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var dto user.CreateUserDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	user, err2 := h.service.SignUp(context.TODO(), dto)
	if err2 != nil {
		return err2
	}

	respWithLogin := fmt.Sprintf("You successfully create user named %s!", user.Login)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respWithLogin))

	return nil
}

// В ответе возвращаем токены в cookie
func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var inp user.SignInUserDTO

	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		return err
	}

	res, err := h.service.SignIn(context.TODO(), inp.Login, inp.Password)
	if err != nil {
		return err
	}

	cookie1 := http.Cookie{
		Name:     "AccessToken",
		Value:    res.AccessToken,
		Path:     "http://localhost:10000/api/",
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   true,
	}

	cookie2 := http.Cookie{
		Name:     "RefreshToken",
		Value:    res.RefreshToken,
		Path:     "http://localhost:10000/api/",
		MaxAge:   2592000,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &cookie2)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user is signed in"))

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

		err = h.service.Verify(context.TODO(), cookie.Value)
		if err != nil {
			h.logger.Errorf("access token is wrong: %v", err)
			return apperror.ErrUnauthorized
		}

		return protectedHandler(w, r)
	}
}

// В ответе возвращаем токены в cookie
func (h *handler) UserRefresh(w http.ResponseWriter, r *http.Request) error {
	cookie1, err := r.Cookie("AccessToken")
	if err != nil {
		h.logger.Errorf("cookie error: %v", err)
		return apperror.ErrUnauthorized
	}

	cookie2, err := r.Cookie("RefreshToken")
	if err != nil {
		h.logger.Errorf("cookie error: %v", err)
		return apperror.ErrUnauthorized
	}

	if cookie2.Value == "" {
		h.logger.Errorf("access token is empty")
		return apperror.ErrUnauthorized
	}

	dto := user.RefreshUserDTO{
		RefreshToken: cookie2.Value,
	}

	res, err := h.service.UserRefresh(context.TODO(), dto)
	if err != nil {
		return err
	}

	cookie1.Value = res.AccessToken
	cookie1.MaxAge = 7200 // 2 hours

	cookie2.Value = res.RefreshToken
	cookie2.MaxAge = 2592000 // 30 days

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user is refreshed"))

	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {

	userID := r.URL.Query().Get("id")

	if userID == "" {
		return apperror.NewAppError(nil, "query param is empty", "param userID is empty", "US-000015")
	}

	user, err := h.service.GetUserByUUID(context.TODO(), userID)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(userBytes)

	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	var dto user.UpdateUserDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.UpdateUser(context.TODO(), dto)
	if err != nil {
		return err
	}

	w.WriteHeader(204)
	w.Write([]byte("user is updated"))

	return nil
}

// TODO
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	var dto user.PartUpdateUserDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	err = h.service.PartUpdateUser(context.TODO(), dto)
	if err != nil {
		return err
	}

	w.WriteHeader(204)
	w.Write([]byte("user is updated"))

	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	var deletingUser user.DeleteUserDTO

	err := json.NewDecoder(r.Body).Decode(&deletingUser)
	if err != nil {
		return err
	}

	err = h.service.DeleteAccount(context.TODO(), deletingUser)
	if err != nil {
		return apperror.NewAppError(err, "wrong id", err.Error(), "US-000009")
	}

	w.WriteHeader(200)
	w.Write([]byte("user is deleted"))

	return nil
}

func (h *handler) CreateEvent(w http.ResponseWriter, r *http.Request) error {
	var dto event.CreateEventDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	eventID, err2 := h.eventService.NewEvent(context.TODO(), dto)
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

func (h *handler) CreateBar(w http.ResponseWriter, r *http.Request) error {
	var dto bar.CreateBarDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	barID, menu, err2 := h.barService.OpenBar(context.TODO(), dto)
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

func (h *handler) NewMenu(w http.ResponseWriter, r *http.Request) error {
	var inp menu.CreateMenuDTO

	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		return err
	}

	menu, err := h.service.CreateMenu(context.TODO(), inp)
	if err != nil {
		return err
	}

	menuBytes, err := json.Marshal(menu)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(menuBytes)

	return nil
}

func (h *handler) AddDrink(w http.ResponseWriter, r *http.Request) error {
	var inp menu.AddDrinkDTO

	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		return err
	}

	menu, err := h.service.AddDrink(context.TODO(), inp)
	if err != nil {
		return err
	}

	menuBytes, err := json.Marshal(menu)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(menuBytes)

	return nil
}
