package user

import (
	"context"
	"time"

	"restapi/internal/domain/session"
	"restapi/pkg/auth"
	"restapi/pkg/hash"
	"restapi/pkg/logging"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Service interface {
	SignUp(ctx context.Context, dto CreateUserDTO) (user User, err error)
	SignIn(ctx context.Context, login, password string) (Tokens, error)
	Verify(ctx context.Context, code string) error
	UserRefresh(ctx context.Context, dto RefreshUserDTO) (Tokens, error)
	DeleteAccount(ctx context.Context, deleteDTO DeleteUserDTO) error
	CreateSession(ctx context.Context, userID string) (Tokens, error)
	UpdateSession(ctx context.Context, userID string) (Tokens, error)
	CreateEvent(ctx context.Context, user User)
}

type service struct {
	repository   Repository
	sessionRepo  session.Repository
	logger       *logging.Logger
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(repository Repository, sessionRepo session.Repository, logger *logging.Logger,
	hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTTL, refreshTTL time.Duration) Service {
	return &service{
		repository:      repository,
		sessionRepo:     sessionRepo,
		logger:          logger,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// TODO for next lesson
func (s *service) SignUp(ctx context.Context, dto CreateUserDTO) (user User, err error) {
	s.logger.Infof("creating user %s", dto.Login)

	user, err = s.repository.Create(ctx, dto)
	if err != nil {
		return
	}

	s.logger.Infof("user %s is created", user.Login)

	return
}

func (s *service) SignIn(ctx context.Context, login, password string) (Tokens, error) {
	// Дает доступ к методам handler, таким как CreateEvent

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repository.GetByCredentials(ctx, login, passwordHash)
	if err != nil {
		// TODO err.UserNotFound
		return Tokens{}, err
	}

	return s.CreateSession(ctx, user.ID)

}

func (s *service) Verify(ctx context.Context, code string) error {

	_, err := s.tokenManager.Parse(code)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) UserRefresh(ctx context.Context, dto RefreshUserDTO) (Tokens, error) {

	userID, err := s.sessionRepo.GetByRefreshToken(ctx, dto.RefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.UpdateSession(ctx, userID)
}

func (s *service) DeleteAccount(ctx context.Context, deleteDTO DeleteUserDTO) error {
	s.logger.Infof("deleting user %s", deleteDTO.ID)

	err := s.repository.Delete(ctx, deleteDTO.ID)
	if err != nil {
		return err
	}

	// TODO удаление всех связанных с юзером сущностей: ивенты --> бары --> меню, сессии, отчеты

	s.logger.Infof("user %s is deleted", deleteDTO.ID)

	return nil
}

// Создаем текущую юзер-сессию, сохраняем ее в бд, возвращаем токены текущей сессии
func (s *service) CreateSession(ctx context.Context, userID string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userID, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := session.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.sessionRepo.SetSession(ctx, userID, session)

	return res, err
}

func (s *service) UpdateSession(ctx context.Context, userID string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userID, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := session.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.sessionRepo.UpdateSession(ctx, userID, session)

	return res, err
}

func (s *service) CreateEvent(ctx context.Context, user User) {

}
