package ingredients_db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/ingredients"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) AddIngredient(ctx context.Context, dto ingredients.AddIngredientDTO) (string, error) {
	q := `
	INSERT INTO ingredients
    	(user_id, event_id, type, name, unit, volume, cost)
	VALUES
    	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var ingrID string

	row := r.client.QueryRow(ctx, q, dto.UserID, dto.EventID, dto.Type, dto.Name, dto.Unit, dto.Volume, dto.Cost)
	err := row.Scan(&ingrID)
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

	return ingrID, nil
}

func (r *repository) AddIngredients(ctx context.Context, dto ingredients.AddIngredientsDTO) ([]string, error) {
	data := EncodeInsertValue(dto.Ingredients, dto.UserID, dto.EventID)

	q := fmt.Sprintf(`
	INSERT INTO ingredients
    	(user_id, event_id, type, name, unit, volume, cost)
	VALUES
    	%s
	RETURNING id
	`, data)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var ingrIDs []string
	var ID string

	for rows.Next() {
		err := rows.Scan(&ID)
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

		ingrIDs = append(ingrIDs, ID)
	}

	return ingrIDs, nil
}

func (r *repository) DeleteIngredient(ctx context.Context, dto ingredients.DeleteIngredientDTO) error {
	q := `
	DELETE FROM 
		ingredients
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
		err := fmt.Errorf("database deleting error: ingredient not found")
		return err
	}

	return nil
}

func (r *repository) DeleteEventIngredients(ctx context.Context, dto ingredients.DeleteEventIngrDTO) error {
	q := `
	DELETE FROM 
		ingredients
	WHERE
		event_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, dto.EventID)
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

	if ct.String() == "DELETE 0" {
		err := fmt.Errorf("database deleting error: any event ingredient not found")
		return err
	}

	return nil
}

func (r *repository) FindIngredient(ctx context.Context,
	dto ingredients.FindIngredientDTO) (ingredients.Ingredient, error) {

	q := `
	SELECT
		id, user_id, event_id, type, name, unit, volume, cost
	FROM 
		ingredients
	WHERE 
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var ingr ingredients.Ingredient

	err := r.client.QueryRow(ctx, q, dto.ID).Scan(&ingr.ID, &ingr.UserID, &ingr.EventID, &ingr.Type, &ingr.Name,
		&ingr.Unit, &ingr.Volume, &ingr.Cost)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return ingredients.Ingredient{}, newErr
		}
		return ingredients.Ingredient{}, err
	}

	return ingr, nil
}

func (r *repository) FindEventIngredients(ctx context.Context,
	dto ingredients.FindEventIngredientsDTO) ([]ingredients.Ingredient, error) {

	q := `
	SELECT
		id, user_id, event_id, type, name, unit, volume, cost
	FROM 
		ingredients
	WHERE 
    	event_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, dto.EventID)
	if err != nil {
		return nil, err
	}

	var Ingrs []ingredients.Ingredient
	var ingr ingredients.Ingredient

	for rows.Next() {
		err := rows.Scan(&ingr.ID, &ingr.UserID, &ingr.EventID, &ingr.Type, &ingr.Name,
			&ingr.Unit, &ingr.Volume, &ingr.Cost)
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

		Ingrs = append(Ingrs, ingr)
	}

	return Ingrs, nil
}

func (r *repository) UpdateIngredient(ctx context.Context, dto ingredients.UpdateIngredientDTO) (string, error) {
	q := `
	UPDATE ingredients
	SET
		type = $2, name = $3, unit = $4, volume = $5, cost = $6
	WHERE
		id = $1
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var updatedID string

	err := r.client.QueryRow(ctx, q, dto.ID, dto.Type, dto.Name, dto.Unit, dto.Volume, dto.Cost).Scan(&updatedID)
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

func NewRepository(client postgresql.Client, logger *logging.Logger) ingredients.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func EncodeInsertValue(ingredients []ingredients.IngredientDataDTO, userID, eventID string) string {
	var buffer bytes.Buffer

	for _, ingr := range ingredients {

		buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', %d, %d)",
			userID, eventID, ingr.Type, ingr.Name, ingr.Unit, ingr.Volume, ingr.Cost))

		// Если это последний ингредиент, то "," в конце не ставим
		if ingr == ingredients[len(ingredients)-1] {
			break
		}

		buffer.WriteString(",")
	}

	return buffer.String()
}
