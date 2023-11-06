package bar_db

import (
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/bar"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	statusOpened = "Opened"
	statusClosed = "Closed"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) CreateBar(ctx context.Context, dto bar.CreateBarDTO) (uint32, error) {
	q := `
	INSERT INTO bars
    	(event_id, name, description, status)
	VALUES
    	($1, $2, $3, $4)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var barID uint32

	row := r.client.QueryRow(ctx, q, dto.EventID, dto.Name, dto.Description, statusOpened)
	err := row.Scan(&barID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return 0, newErr
		}

		return 0, err
	}

	return barID, nil
}

// Не удаляет бар из таблицы, а устанавливает статус "Closed"
func (r *repository) CloseBar(ctx context.Context, dto bar.CloseBarDTO) error {
	q := `
	UPDATE bars
	SET
    	status = $2
	WHERE
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, dto.ID, statusClosed)
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
		err := fmt.Errorf("database updating error: bar not found")
		return err
	}

	return nil
}

func (r *repository) GetOrders(ctx context.Context, dto bar.GetOrdersDTO) ([]string, error) {
	q := `
	SELECT
    	orders
	FROM 
    	bars
	WHERE 
    	event_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, dto.EventID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}

		return nil, err
	}

	ordersID := make([]string, 0)

	for rows.Next() {
		var ordrs []string

		err = rows.Scan(&ordrs)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return nil, newErr
			}

			return nil, err
		}

		ordersID = append(ordersID, ordrs...)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ordersID, nil
}

func (r *repository) GetBarOrders(ctx context.Context, dto bar.GetBarOrdersDTO) ([]string, error) {
	q := `
	SELECT
    	orders
	FROM 
    	bars
	WHERE 
    	id = $1 AND event_id = $2
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ordersID := make([]string, 0)

	err := r.client.QueryRow(ctx, q, dto.ID, dto.EventID).Scan(&ordersID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}

		return nil, err
	}

	return ordersID, nil
}

func (r *repository) UpdateInfo(ctx context.Context, dto bar.UpdateBarDTO) error {
	q := `
	UPDATE bars
	SET
		name = $2, description = $2, orders = $3, session_url = $4
	WHERE
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, dto.ID, dto.Name, dto.Description, dto.Orders, dto.SessionURL)
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
		err := fmt.Errorf("database updating error: bar not found")
		return err
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) bar.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
