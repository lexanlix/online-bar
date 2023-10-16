package menu_db

import (
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/menu"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) CreateMenu(ctx context.Context, dto menu.CreateMenuDTO) (string, error) {
	q := `
	INSERT INTO menu
		(user_id, name, drinks, total_cost)
	VALUES
    	($1, $2, $3, $4)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var menuID string

	drinksString := GetInsertValue(dto.Drinks)

	row := r.client.QueryRow(ctx, q, dto.UserID, dto.Name, drinksString, 600)
	err := row.Scan(&menuID)
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

	return menuID, nil
}

// TODO
func (r *repository) DeleteMenu(ctx context.Context, dto menu.DeleteMenuDTO) (string, error) {
	q := `
	UPDATE bars
	SET
    	status = $2
	WHERE
    	id = $1
	RETURNING status
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var status string
	/*
		err := r.client.QueryRow(ctx, q, dto.ID, statusClosed).Scan(&status)
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
	*/
	return status, nil
}

// TODO
func (r *repository) FindMenu(ctx context.Context, dto menu.FindMenuDTO) (menu.Menu, error) {
	q := `
	SELECT
    	orders
	FROM 
    	bars
	WHERE 
    	id = $1 AND event_id = $2
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var mn menu.Menu
	/*
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
	*/
	return mn, nil
}

// TODO
func (r *repository) UpdateMenu(ctx context.Context, dto menu.UpdateMenuDTO) (string, error) {
	q := `
	UPDATE bars
	SET
		name = $2, description = $2, orders = $3, session_url = $4
	WHERE
		id = $1
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var updatedID string
	/*
		err := r.client.QueryRow(ctx, q, dto.ID, dto.Name, dto.Description, dto.Orders, dto.SessionURL).Scan(&updatedID)
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
	*/
	return updatedID, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) menu.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func GetInsertValue(drinks map[string][]menu.Drink) string {
	var insertString string

	return insertString
}
