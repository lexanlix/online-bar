package session

import (
	"context"
)

type Repository interface {
	SetSession(ctx context.Context, userID string, session Session) error
	GetByRefreshToken(ctx context.Context, token string) (string, error)
	UpdateSession(ctx context.Context, userID string, session Session) error
}
