package user_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/adapters"
	"restapi/internal/apperror"
	"restapi/internal/domain/user"

	"restapi/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// Подсказка, что структура реализует интерфейс
var _ adapters.Handler = &handler{}

const (
	usersURL        = "/api/users"
	userURL         = "/api/users/:uuid"
	registerUserURL = "/api/register"
	signInURL       = "/api/signIn"
	refreshURL      = "/api/auth/refresh"
	updateUserURL   = "/api/update"
	pUpdateUserURL  = "/api/update/part"
	deleteUserURL   = "/api/user/delete"
)

type signInInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type handler struct {
	service user.Service
	logger  *logging.Logger
}

func NewHandler(logger *logging.Logger, service user.Service) adapters.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func Protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Protected!\n")
}

func (h *handler) Register(router *httprouter.Router) {

	router.HandlerFunc(http.MethodPost, registerUserURL, apperror.Middleware(h.UserSignUp))
	router.HandlerFunc(http.MethodPost, signInURL, apperror.Middleware(h.SignIn))
	router.HandlerFunc(http.MethodPost, refreshURL, apperror.Middleware(h.UserRefresh))

	// Обработчики, доступные пользователям, вошедшим в аккаунт (которые имеют AccessToken)
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.Verify(h.GetUserByUUID)))
	router.HandlerFunc(http.MethodPut, updateUserURL, apperror.Middleware(h.Verify(h.UpdateUser)))
	router.HandlerFunc(http.MethodPatch, pUpdateUserURL, apperror.Middleware(h.Verify(h.PartiallyUpdateUser)))
	router.HandlerFunc(http.MethodDelete, deleteUserURL, apperror.Middleware(h.Verify(h.DeleteUser)))
}

func (h *handler) UserSignUp(w http.ResponseWriter, r *http.Request) error {
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

// В ответе возвращаем токены в json
func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var inp signInInput

	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		return err
	}

	res, err := h.service.SignIn(context.TODO(), inp.Login, inp.Password)
	if err != nil {
		return err
	}

	tokenResponse := tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	respBytes, err := json.Marshal(tokenResponse)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) Verify(protectedHandler apperror.AppHandler) apperror.AppHandler {

	return func(w http.ResponseWriter, r *http.Request) error {
		code := r.Header.Get("Bearer")

		if code == "" {
			h.logger.Errorf("access token is empty")
			return apperror.ErrUnauthorized
		}

		err := h.service.Verify(context.TODO(), code)
		if err != nil {
			h.logger.Errorf("access token is wrong: %v", err)
			return apperror.ErrUnauthorized
		}

		return protectedHandler(w, r)
	}
}

// В ответе возвращаем токены в json
func (h *handler) UserRefresh(w http.ResponseWriter, r *http.Request) error {
	var dto user.RefreshUserDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return err
	}

	res, err := h.service.UserRefresh(context.TODO(), dto)
	if err != nil {
		return err
	}

	tokenResponse := tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	respBytes, err := json.Marshal(tokenResponse)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)

	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	// w.WriteHeader(200)
	// w.Write([]byte("this is user by uuid"))

	// return nil

	return apperror.NewAppError(nil, "test", "test", "t123")
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is update user"))

	return nil
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is partially update user"))

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
