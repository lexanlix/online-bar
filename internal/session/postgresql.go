package session

import (
	"context"
	"fmt"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"
	"time"
)

type SessionRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (sr *SessionRepository) SetSession(ctx context.Context, userID string, session Session) error {
	q := `	
	INSERT INTO sessions 
		(user_id, refresh_token, expires_at, last_visit_at)
	VALUES 
		($1, $2, $3, $4)
	` // CHECK REQUEST!!! //

	sr.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	row := sr.client.QueryRow(ctx, q, userID, session.RefreshToken, session.ExpiresAt, time.Now())

	err := row.Scan()

	if err.Error() != "no rows in result set" {
		return err
	}

	return nil
}

func NewSessionRepository(client postgresql.Client, logger *logging.Logger) Repository {
	return &SessionRepository{
		client: client,
		logger: logger,
	}
}
