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
	SignUp(ctx context.Context, dto CreateUserDTO) (User, error)
	SignIn(ctx context.Context, login, password string) (Tokens, error)
	Verify(ctx context.Context, code string) error
	UserRefresh(ctx context.Context, dto RefreshUserDTO) (Tokens, error)
	UpdateUser(ctx context.Context, dto UpdateUserDTO) error
	PartUpdateUser(ctx context.Context, dto PartUpdateUserDTO) error
	GetUserByUUID(ctx context.Context, userID string) (User, error)
	DeleteAccount(ctx context.Context, deleteDTO DeleteUserDTO) error
	CreateSession(ctx context.Context, userID string) (Tokens, error)
	UpdateSession(ctx context.Context, userID string) (Tokens, error)
	CreateEvent(ctx context.Context, user User)

	//CreateMenu(ctx context.Context, dto menu.CreateMenuDTO) (menu.Menu, error)
	//AddDrink(ctx context.Context, dto menu.Drink) (menu.Menu, error)
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
	s.logger.Infof("refreshing user")

	userID, err := s.sessionRepo.GetByRefreshToken(ctx, dto.RefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.UpdateSession(ctx, userID)
}

func (s *service) UpdateUser(ctx context.Context, dto UpdateUserDTO) error {
	s.logger.Infof("updating user %s", dto.ID)

	passwordHash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return err
	}

	updateUser := User{
		ID:           dto.ID,
		Name:         dto.Name,
		Login:        dto.Login,
		PasswordHash: passwordHash,
		OneTimeCode:  dto.OneTimeCode,
	}

	err = s.repository.Update(ctx, updateUser)
	if err != nil {
		return err
	}

	s.logger.Infof("updating user %s", dto.ID)
	return nil
}

func (s *service) PartUpdateUser(ctx context.Context, dto PartUpdateUserDTO) error {
	s.logger.Infof("updating user '%s' field '%s'", dto.ID, dto.Key)

	var err error

	if dto.Key == "password" {
		dto.Key = "password_hash"
		dto.Value, err = s.hasher.Hash(dto.Value)
		if err != nil {
			return err
		}
	}

	err = s.repository.PartUpdate(ctx, dto)
	if err != nil {
		return err
	}

	s.logger.Infof("user '%s' field '%s' is updated", dto.ID, dto.Key)
	return nil
}

func (s *service) GetUserByUUID(ctx context.Context, userID string) (User, error) {
	s.logger.Infof("getting user %s", userID)

	user, err := s.repository.GetByUUID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	s.logger.Infof("user %s is received", user.Name)
	return user, nil
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

// TODO
func (s *service) CreateEvent(ctx context.Context, user User) {

}

// // Должно выполняться на клиенте пользователя
// func (s *service) CreateMenu(ctx context.Context, dto menu.CreateMenuDTO) (menu.Menu, error) {
// 	newMenu := menu.NewMenu(dto.ID, dto.Name, dto.Drinks)
// 	newMenu.UpdateTotalCost()

// 	return newMenu, nil
// }

// // Должно выполняться на клиенте пользователя
// // Добавляет напиток в меню (в указанную группу) и обновляет total_cost
// func (s *service) AddDrink(ctx context.Context, dto menu.Drink) (menu.Menu, error) {
// 	drinks := make(map[string][]menu.Drink, 0)

// 	emptyMeny := menu.Menu{
// 		Drinks: drinks,
// 	}

// 	err := emptyMeny.AddDrink(dto)
// 	if err != nil {
// 		return menu.Menu{}, err
// 	}

// 	return emptyMeny, nil
// }
