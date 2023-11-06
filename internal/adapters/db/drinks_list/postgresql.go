package drinks_list_db

import (
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/drinks_list"
	"restapi/internal/domain/menu"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	noCat = iota
	Beer
	Cider
	LongDrink
	NonAlco
	ShortDrink
	ShotDrink
	StrongAlco
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) AddUserDrink(ctx context.Context, dto drinks_list.AddUserDrinkDTO) (string, error) {
	q := `
	INSERT INTO user_drinks
		(user_id, name, category, cooking_method, composition, ice_type, price, bars_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var drinkID string

	err := r.client.QueryRow(ctx, q, dto.UserID, dto.Name, dto.Category, dto.Cooking_method, dto.Composition,
		dto.OrderIceType, dto.Price, dto.BarsID).Scan(&drinkID)
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

	return drinkID, nil
}

func (r *repository) DeleteUserDrink(ctx context.Context, dto drinks_list.DeleteUserDrinkDTO) error {
	q := `
	DELETE FROM user_drinks
	WHERE
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, dto.ID)
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

	if ct.String() != "DELETE 1" {
		err := fmt.Errorf("database deleting error: user drink not found")
		return err
	}

	return nil
}

func (r *repository) FindUserDrink(ctx context.Context, dto drinks_list.FindUserDrinkDTO) (menu.Drink, error) {
	q := `
	SELECT
		id, name, category, cooking_method, composition, ice_type, price, bars_id
	FROM 
		user_drinks
	WHERE 
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var dr menu.Drink

	err := r.client.QueryRow(ctx, q, dto.ID).Scan(&dr.ID, &dr.Name, &dr.Category, &dr.Cooking_method, &dr.Composition,
		&dr.OrderIceType, &dr.Price, &dr.BarsID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return menu.Drink{}, newErr
		}
		return menu.Drink{}, err
	}

	return dr, nil
}

func (r *repository) FindUserDrinks(ctx context.Context,
	dto drinks_list.FindUserDrinksDTO) (drinks_list.RespFindUDrinks, error) {
	q := `
	SELECT
		id, name, category, cooking_method, composition, ice_type, price, bars_id
	FROM 
		user_drinks
	WHERE
    	user_id = $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, dto.UserID)
	if err != nil {
		return drinks_list.RespFindUDrinks{}, err
	}

	var UserDrinks drinks_list.RespFindUDrinks

	for rows.Next() {
		var dr menu.Drink

		err := rows.Scan(&dr.ID, &dr.Name, &dr.Category, &dr.Cooking_method, &dr.Composition,
			&dr.OrderIceType, &dr.Price, &dr.BarsID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return drinks_list.RespFindUDrinks{}, newErr
			}

			return drinks_list.RespFindUDrinks{}, err
		}

		UserDrinks.Drinks = append(UserDrinks.Drinks, dr)
	}

	return UserDrinks, nil
}

func (r *repository) UpdateUserDrink(ctx context.Context, dto drinks_list.UpdateUserDrinkDTO) (string, error) {
	q := `
	UPDATE user_drinks
	SET
		name = $2, category = $3, cooking_method = $4, composition = $5, ice_type = $6, price = $7, bars_id = $8
	WHERE
		id = $1
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var updatedID string

	err := r.client.QueryRow(ctx, q, dto.ID, dto.Name, dto.Category, dto.Cooking_method, dto.Composition, dto.OrderIceType,
		dto.Price, dto.BarsID).Scan(&updatedID)
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

func NewRepository(client postgresql.Client, logger *logging.Logger) drinks_list.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
