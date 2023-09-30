package session

import "context"

type Repository interface {
	SetSession(ctx context.Context, userID string, session Session) error
}
