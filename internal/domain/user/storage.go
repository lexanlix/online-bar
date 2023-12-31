package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, dto CreateUserDTO) (user User, err error)
	GetByCredentials(ctx context.Context, login, passwordHash string) (User, error)
	GetByUUID(ctx context.Context, userID string) (User, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	PartUpdate(ctx context.Context, dto PartUpdateUserDTO) error
	Delete(ctx context.Context, id string) error
}
