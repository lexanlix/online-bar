package event_db

import (
	"context"
	"fmt"
	"restapi/internal/domain/event"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"
)

const (
	statusCreated   = "Created"
	statusCompleted = "Completed"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) CreateEvent(ctx context.Context, dto event.CreateEventDTO) (string, error) {
	q := `
	INSERT INTO events
    	(host_id, name, description, participants_number, date_time, status)
	VALUES
    	($1, $2, $3, $4, $5, $6)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var eventID string

	row := r.client.QueryRow(ctx, q, dto.HostID, dto.Name, dto.Description, dto.ParticipantsNumber,
		dto.DateTime, statusCreated)
	err := row.Scan(&eventID)
	if err != nil {
		return "", err
	}

	return eventID, nil
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
		return "", err
	}

	return status, nil
}

func (r *repository) FindAllUserEvents(ctx context.Context, dto event.FindAllEventsDTO) ([]event.Event, error) {
	q := `
	SELECT 
    	id, host_id, name, description, participants_number, date_time
	FROM 
    	events
	WHERE
    	host_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, dto.HostID)
	if err != nil {
		return nil, err
	}

	events := make([]event.Event, 0)

	for rows.Next() {
		var evnt event.Event

		err = rows.Scan(&evnt.ID, &evnt.HostID, &evnt.Name, &evnt.Description,
			&evnt.ParticipantsNumber, &evnt.DateTime)
		if err != nil {
			return nil, err
		}

		events = append(events, evnt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *repository) FindOneUserEvent(ctx context.Context, dto event.FindEventDTO) (event.Event, error) {
	q := `
	SELECT 
    	id, host_id, name, description, participants_number, date_time
	FROM 
    	events
	WHERE
    	id = $1 AND host_id = $2
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var evnt event.Event

	err := r.client.QueryRow(ctx, q, dto.ID, dto.HostID).Scan(&evnt.ID, &evnt.HostID, &evnt.Name,
		&evnt.Description, &evnt.ParticipantsNumber, &evnt.DateTime)
	if err != nil {
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
		return "", err
	}

	return updatedID, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) event.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}