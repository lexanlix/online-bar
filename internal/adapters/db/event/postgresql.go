package event_db

import (
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/event"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	statusCreated   = "Created"
	statusActive    = "Active"
	statusCompleted = "Completed"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) CreateEvent(ctx context.Context, dto event.CreateEventDTO, shopList []string) (string, error) {
	q := `
	INSERT INTO events
    	(user_id, name, description, participants_number, date_time, status, menu_id, shopping_list)
	VALUES
    	($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var eventID string

	row := r.client.QueryRow(ctx, q, dto.UserID, dto.Name, dto.Description, dto.ParticipantsNumber,
		dto.DateTime, statusCreated, dto.MenuID, shopList)
	err := row.Scan(&eventID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	return eventID, nil
}

func (r *repository) SetActive(ctx context.Context, event_id string) (string, error) {
	q := `
	UPDATE events
	SET 
		status = $2
	WHERE 
		id = $1
	RETURNING status
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var status string

	err := r.client.QueryRow(ctx, q, event_id, statusActive).Scan(&status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	return status, nil
}

// Не удаляет ивент из таблицы, а устанавливает статус "Completed"
func (r *repository) DeleteEvent(ctx context.Context, dto event.CompleteEventDTO) (string, error) {
	q := `
	UPDATE events
	SET 
		status = $2
	WHERE 
		id = $1
	RETURNING status
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var status string

	err := r.client.QueryRow(ctx, q, dto.ID, statusCompleted).Scan(&status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	return status, nil
}

func (r *repository) FindAllUserEvents(ctx context.Context, dto event.FindAllEventsDTO) (event.RespAllEvents, error) {
	q := `
	SELECT 
    	id, user_id, name, description, participants_number, date_time, status
	FROM 
    	events
	WHERE
		user_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, dto.UserID)
	if err != nil {
		return event.RespAllEvents{}, err
	}

	var resp event.RespAllEvents

	for rows.Next() {
		var evnt event.Event

		err = rows.Scan(&evnt.ID, &evnt.UserID, &evnt.Name, &evnt.Description,
			&evnt.ParticipantsNumber, &evnt.DateTime, &evnt.Status)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return event.RespAllEvents{}, newErr
			}

			return event.RespAllEvents{}, err
		}

		resp.Events = append(resp.Events, evnt)
	}

	if err = rows.Err(); err != nil {
		return event.RespAllEvents{}, err
	}

	return resp, nil
}

func (r *repository) FindUserEvent(ctx context.Context, dto event.FindEventDTO) (event.Event, error) {
	q := `
	SELECT 
    	id, user_id, name, description, participants_number, date_time, status, menu_id, shopping_list
	FROM 
    	events
	WHERE
    	id = $1 AND user_id = $2
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var evnt event.Event

	err := r.client.QueryRow(ctx, q, dto.ID, dto.UserID).Scan(&evnt.ID, &evnt.UserID, &evnt.Name, &evnt.Description,
		&evnt.ParticipantsNumber, &evnt.DateTime, &evnt.Status, &evnt.MenuID, &evnt.ShoppingList)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return event.Event{}, newErr
		}

		return event.Event{}, err
	}

	return evnt, nil
}

func (r *repository) UpdateEvent(ctx context.Context, dto event.UpdateEventDTO) (string, error) {
	q := `
	UPDATE events
	SET 
		name = $2, description = $3, participants_number = $4, date_time = $5
	WHERE 
		id = $1
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var updatedID string

	err := r.client.QueryRow(ctx, q, dto.ID, dto.Name, dto.Description, dto.ParticipantsNumber,
		dto.DateTime).Scan(&updatedID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	return updatedID, nil
}

func (r *repository) UpdateIceTypesNum(ctx context.Context, onlyOneIceType bool, eventID string) error {
	q := `
	UPDATE events
	SET 
		only_one_ice_type = $2
	WHERE 
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, eventID, onlyOneIceType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}

		return err
	}

	if ct.String() != "UPDATE 1" {
		err := fmt.Errorf("database deleting error: event not found")
		return err
	}

	return nil
}

// Вызывается при создании сессий баров. Проверить!
func (r *repository) GetIceTypesNum(ctx context.Context, eventID string) (bool, error) {
	q := `
	SELECT 
		only_one_ice_type
	FROM
		events
	WHERE 
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var onlyOneIceType bool

	err := r.client.QueryRow(ctx, q, eventID).Scan(&onlyOneIceType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return false, newErr
		}

		return false, err
	}

	return onlyOneIceType, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) event.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
