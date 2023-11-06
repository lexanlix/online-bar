package session_db

import (
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/session"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type SessionRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (sr *SessionRepository) SetSession(ctx context.Context, userID string, session session.Session) error {
	q := `	
	INSERT INTO sessions 
		(user_id, refresh_token, expires_at, last_visit_at)
	VALUES 
		($1, $2, $3, $4)
	`

	sr.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	row := sr.client.QueryRow(ctx, q, userID, session.RefreshToken, session.ExpiresAt, time.Now())

	err := row.Scan()

	if err.Error() != "no rows in result set" {
		return err
	}

	return nil
}

func (sr *SessionRepository) UpdateSession(ctx context.Context, userID string, session session.Session) error {
	q := `	
	UPDATE 
		sessions
	SET
    	refresh_token = $1,
    	expires_at = $2,
    	last_visit_at = $3
	WHERE
    	user_id = $4
	`
	sr.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := sr.client.Exec(ctx, q, session.RefreshToken, session.ExpiresAt, time.Now(), userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			sr.logger.Error(newErr)
			return newErr
		}

		return err
	}

	if ct.String() != "UPDATE 1" {
		err := fmt.Errorf("database updating error: user session not found")
		return err
	}

	return nil
}

func (sr *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (string, error) {
	q := `
	SELECT 
		user_id
	FROM
		public.sessions
	WHERE 
		(refresh_token = $1) AND (expires_at > $2)
	`
	sr.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	row := sr.client.QueryRow(ctx, q, token, time.Now())

	var userID string
	err := row.Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func NewSessionRepository(client postgresql.Client, logger *logging.Logger) session.Repository {
	return &SessionRepository{
		client: client,
		logger: logger,
	}
}
