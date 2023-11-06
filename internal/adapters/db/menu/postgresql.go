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

func (r *repository) CreateMenu(ctx context.Context, dto menu.MenuDTO, totalCost uint32) (string, error) {
	tx, err := r.client.Begin(ctx)
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

	q_menu := `
	INSERT INTO menu
		(user_id, name, total_cost)
	VALUES
    	($1, $2, $3)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q_menu)))

	var menuID string

	row := tx.QueryRow(ctx, q_menu, dto.UserID, dto.Name, totalCost)
	err = row.Scan(&menuID)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	q_drinks := `
	INSERT INTO menu_drinks
		(menu_id, name, category, cooking_method, composition, ice_type, price, bars_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, drinks := range dto.Drinks {
		for _, drink := range drinks {
			r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q_menu)))

			_, err := tx.Exec(ctx, q_drinks, menuID, drink.Name, drink.Category, drink.Cooking_method, drink.Composition,
				drink.OrderIceType, drink.Price, drink.BarsID)
			if err != nil {
				tx.Rollback(ctx)
				tx.Conn().Close(ctx)

				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					pgErr = err.(*pgconn.PgError)
					newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
					r.logger.Error(newErr)
					return "", newErr
				}

				return "", err
			}
		}
	}

	tx.Commit(ctx)
	tx.Conn().Close(ctx)

	return menuID, nil
}

func (r *repository) DeleteMenu(ctx context.Context, dto menu.DeleteMenuDTO) error {
	tx, err := r.client.Begin(ctx)
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

	q := `
	DELETE FROM 
		menu
	WHERE
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := tx.Exec(ctx, q, dto.ID)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

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
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		err := fmt.Errorf("database deleting error: menu not found")
		return err
	}

	q = `
	DELETE FROM 
		menu_drinks
	WHERE
    	menu_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	_, err = tx.Exec(ctx, q, dto.ID)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}

		return err
	}

	tx.Commit(ctx)
	tx.Conn().Close(ctx)

	return nil
}

func (r *repository) FindMenu(ctx context.Context, dto menu.FindMenuDTO) (menu.Menu, error) {
	q := `
	SELECT
		id, user_id, name, total_cost
	FROM 
    	menu
	WHERE 
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var mn menu.Menu

	err := r.client.QueryRow(ctx, q, dto.ID).Scan(&mn.ID, &mn.UserID, &mn.Name, &mn.TotalCost)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return menu.Menu{}, newErr
		}
		return menu.Menu{}, err
	}

	q = `
	SELECT
		category, id, name, cooking_method, composition, ice_type, price, bars_id
	FROM 
		menu_drinks
	WHERE 
		menu_id = $1
	GROUP BY id
	ORDER BY category ASC;
	`

	rows, err := r.client.Query(ctx, q, dto.ID)
	if err != nil {
		return menu.Menu{}, err
	}

	Drinks := make(map[string][]menu.Drink, 0)

	for rows.Next() {
		var drink menu.Drink

		err := rows.Scan(&drink.Category, &drink.ID, &drink.Name, &drink.Cooking_method,
			&drink.Composition, &drink.OrderIceType, &drink.Price, &drink.BarsID)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return menu.Menu{}, newErr
			}
			return menu.Menu{}, err
		}

		Drinks[drink.Category] = append(Drinks[drink.Category], drink)
	}

	mn.Drinks = Drinks

	return mn, nil
}

func (r *repository) FindUserMenus(ctx context.Context, dto menu.UserMenusDTO) (menu.RespUserMenus, error) {
	q := `
	SELECT
    	id, name
	FROM 
		menu
	WHERE
    	user_id = $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, dto.UserID)
	if err != nil {
		return menu.RespUserMenus{}, err
	}

	var UserMenus menu.RespUserMenus
	var mn menu.UserMenu

	for rows.Next() {
		err := rows.Scan(&mn.ID, &mn.Name)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return menu.RespUserMenus{}, newErr
			}

			return menu.RespUserMenus{}, err
		}

		UserMenus.Menus = append(UserMenus.Menus, mn)
	}

	return UserMenus, nil
}

func (r *repository) UpdateMenu(ctx context.Context, dto menu.UpdateMenuDTO, totalCost uint32) (string, error) {
	tx, err := r.client.Begin(ctx)

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

	q := `
	UPDATE menu
	SET
		name = $2, total_cost = $3
	WHERE
		id = $1
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var updatedID string

	err = tx.QueryRow(ctx, q, dto.ID, dto.Name, totalCost).Scan(&updatedID)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	for _, drinks := range dto.Drinks {
		for _, drink := range drinks {
			q = `
			UPDATE menu_drinks
			SET
				name = $2, category = $3, cooking_method = $4, composition = $5, 
				ice_type = $6, price = $7, bars_id = $8
			WHERE
				id = $1
			`
			r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

			ct, err := tx.Exec(ctx, q, drink.ID, drink.Name, drink.Category, drink.Cooking_method, drink.Composition,
				drink.OrderIceType, drink.Price, drink.BarsID)

			if err != nil {
				tx.Rollback(ctx)
				tx.Conn().Close(ctx)

				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					pgErr = err.(*pgconn.PgError)
					newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
					r.logger.Error(newErr)
					return "", newErr
				}

				return "", err
			}

			if ct.String() != "UPDATE 1" {
				tx.Rollback(ctx)
				tx.Conn().Close(ctx)

				err := fmt.Errorf("database updating error: drink %s not found", drink.ID)
				return "", err
			}
		}
	}

	tx.Commit(ctx)
	tx.Conn().Close(ctx)

	return updatedID, nil
}

func (r *repository) UpdateNameMenu(ctx context.Context, dto menu.UpdateMenuNameDTO) error {
	q := `
	UPDATE menu
	SET
		name = $2
	WHERE
		id = $1
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	err := r.client.QueryRow(ctx, q, dto.ID, dto.Name).Scan(&dto.ID)
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

	return nil
}

func (r *repository) AddDrink(ctx context.Context, dto menu.AddDrinkDTO) (string, error) {
	tx, err := r.client.Begin(ctx)
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

	q := `
	INSERT INTO menu_drinks
		(menu_id, name, category, cooking_method, composition, ice_type, price, bars_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING
    	id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var drinkID string

	err = tx.QueryRow(ctx, q, dto.MenuID, dto.Drink.Name, dto.Drink.Category, dto.Drink.Cooking_method,
		dto.Drink.Composition, dto.Drink.OrderIceType, dto.Drink.Price, dto.Drink.BarsID).Scan(&drinkID)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	q = `
	UPDATE menu
	SET
	    total_cost = total_cost + $2
	WHERE
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := tx.Exec(ctx, q, dto.MenuID, dto.Drink.Price)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}

		return "", err
	}

	if ct.String() != "UPDATE 1" {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		err := fmt.Errorf("database updating error: menu not found")
		return "", err
	}

	tx.Commit(ctx)
	tx.Conn().Close(ctx)

	return drinkID, nil
}

func (r *repository) DeleteDrink(ctx context.Context, dto menu.DeleteDrinkDTO) error {
	tx, err := r.client.Begin(ctx)
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

	q := `
	DELETE FROM menu_drinks
	WHERE
		id = $1
	RETURNING
		menu_id, price
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var menuID string
	var price int

	err = tx.QueryRow(ctx, q, dto.DrinkID).Scan(&menuID, &price)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}

		return err
	}

	q = `
	UPDATE menu
	SET
	    total_cost = total_cost - $2
	WHERE
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := tx.Exec(ctx, q, menuID, price)
	if err != nil {
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

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
		tx.Rollback(ctx)
		tx.Conn().Close(ctx)

		err := fmt.Errorf("database updating error: menu not found")
		return err
	}

	tx.Commit(ctx)
	tx.Conn().Close(ctx)

	return nil
}

func (r *repository) FindUserDrink(ctx context.Context, drID string) (menu.NewDrinkDTO, error) {
	q := `
	SELECT
		name, category, cooking_method, composition, ice_type, price, bars_id
	FROM 
		user_drinks
	WHERE
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var newDr menu.NewDrinkDTO

	row := r.client.QueryRow(ctx, q, drID)

	err := row.Scan(&newDr.Name, &newDr.Category, &newDr.Cooking_method, &newDr.Composition, &newDr.OrderIceType,
		&newDr.Price, &newDr.BarsID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return menu.NewDrinkDTO{}, newErr
		}

		return menu.NewDrinkDTO{}, err
	}

	return newDr, nil
}

func (r *repository) FindUserDrinks(ctx context.Context, drIDs []string) ([]menu.Drink, error) {
	var ids string

	for i, id := range drIDs {
		ids += fmt.Sprintf("'%s'", id)

		if i == len(drIDs)-1 {
			break
		}

		ids += ","
	}

	q := fmt.Sprintf(`
	SELECT
		id, name, category, cooking_method, composition, ice_type, price, bars_id
	FROM 
		user_drinks
	WHERE
    	id IN (%s)
	`, ids)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var drinks []menu.Drink

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
				return nil, newErr
			}

			return nil, err
		}

		drinks = append(drinks, dr)
	}

	return drinks, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) menu.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
