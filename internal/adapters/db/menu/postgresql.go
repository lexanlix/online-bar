package menu_db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/menu"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"
	"strings"

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

	sub_q := EncodeInsertValue(&dto.Drinks, menuID)

	q_drinks := fmt.Sprintf(`
	INSERT INTO menu_drinks
		(menu_id, name, category, cooking_method, composition, ice_type, price, bars_id)
	VALUES
    	%s
	RETURNING
    	id
	`, sub_q)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q_menu)))

	_, err = tx.Exec(ctx, q_drinks)
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
		var composition string

		var drink menu.Drink

		err := rows.Scan(&drink.Category, &drink.ID, &drink.Name, &drink.Cooking_method,
			&composition, &drink.OrderIceType, &drink.Price, &drink.BarsID)

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

		drink.Composition, err = DecodeDrinkComposition(composition)
		if err != nil {
			err = fmt.Errorf("decode sql request error: %v", err)
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
			comp := EncodeUpdateValue(&drink.Composition)

			q = fmt.Sprintf(`
			UPDATE menu_drinks
			SET
				name = $2, category = $3, cooking_method = $4, composition = %s, 
				ice_type = $5, price = $6, bars_id = $7
			WHERE
				id = $1
			`, comp)
			r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

			ct, err := tx.Exec(ctx, q, drink.ID, drink.Name, drink.Category, drink.Cooking_method,
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

	comp := EncodeUpdateValue(&dto.Drink.Composition)

	q := fmt.Sprintf(`
	INSERT INTO menu_drinks
		(menu_id, name, category, cooking_method, composition, ice_type, price, bars_id)
	VALUES
		($1, $2, $3, $4, %s, $5, $6, $7)
	RETURNING
    	id
	`, comp)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var drinkID string

	err = tx.QueryRow(ctx, q, dto.MenuID, dto.Drink.Name, dto.Drink.Category, dto.Drink.Cooking_method,
		dto.Drink.OrderIceType, dto.Drink.Price, dto.Drink.BarsID).Scan(&drinkID)
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
	var composition string

	row := r.client.QueryRow(ctx, q, drID)

	err := row.Scan(&newDr.Name, &newDr.Category, &newDr.Cooking_method, &composition, &newDr.OrderIceType,
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

	newDr.Composition, err = DecodeDrinkComposition(composition)
	if err != nil {
		err = fmt.Errorf("decode sql request error: %v", err)
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
    	id IN (%s);
	`, ids)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var drinks []menu.Drink

	for rows.Next() {
		var composition string
		var dr menu.Drink

		err := rows.Scan(&dr.ID, &dr.Name, &dr.Category, &dr.Cooking_method, &composition,
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

		dr.Composition, err = DecodeDrinkComposition(composition)
		if err != nil {
			err = fmt.Errorf("decode sql request error: %v", err)
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

func EncodeInsertValue(drinkGroups *map[string][]menu.NewDrinkDTO, menuID string) string {
	var buffer bytes.Buffer
	var categories []string

	for category := range *drinkGroups {
		categories = append(categories, category)
	}

	for _, category := range categories {
		for _, drink := range (*drinkGroups)[category] {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', CAST(", menuID, drink.Name, drink.Category,
				drink.Cooking_method))

			buffer.WriteString(fmt.Sprintf("(%d, CAST(ARRAY[", drink.Composition.IceBulk))

			for _, liquid := range drink.Composition.Liquids {
				buffer.WriteString(fmt.Sprintf("('%s', '%s', %d)", liquid.Name, liquid.Unit, liquid.Volume))

				// Если это последний элемент массива, то "," в конце не ставим
				if liquid == drink.Composition.Liquids[len(drink.Composition.Liquids)-1] {
					break
				}

				buffer.WriteString(",")
			}
			buffer.WriteString("] AS Liquid []), CAST(ARRAY[")

			for _, solidBulk := range drink.Composition.SolidsBulk {
				buffer.WriteString(fmt.Sprintf("('%s', '%s', %d)", solidBulk.Name, solidBulk.Unit, solidBulk.Volume))

				// Если это последний элемент массива, то "," в конце не ставим
				if solidBulk == drink.Composition.SolidsBulk[len(drink.Composition.SolidsBulk)-1] {
					break
				}

				buffer.WriteString(",")
			}
			buffer.WriteString("] AS Solid_bulk []), CAST(ARRAY[")

			for _, solidUnit := range drink.Composition.SolidsUnit {
				buffer.WriteString(fmt.Sprintf("('%s', %d)", solidUnit.Name, solidUnit.Volume))

				// Если это последний элемент массива, то "," в конце не ставим
				if solidUnit == drink.Composition.SolidsUnit[len(drink.Composition.SolidsUnit)-1] {
					break
				}

				buffer.WriteString(",")
			}

			buffer.WriteString(fmt.Sprintf("] AS Solid_unit [])) AS Composition), '%s', %d, ARRAY[",
				drink.OrderIceType, drink.Price))

			for _, barId := range drink.BarsID {
				buffer.WriteString(fmt.Sprintf("%d", barId))

				// Если это последний элемент массива, то "," в конце не ставим
				if barId == drink.BarsID[len(drink.BarsID)-1] {
					break
				}

				buffer.WriteString(",")
			}

			buffer.WriteString("])")

			// Если это последний напиток последней категории, то "," в конце не ставим
			if drink.Name == (*drinkGroups)[category][len((*drinkGroups)[category])-1].Name && drink.Category == categories[len(categories)-1] {
				break
			}

			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

func EncodeUpdateValue(comp *menu.Composition) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("CAST((%d, CAST(ARRAY[", comp.IceBulk))

	for _, liquid := range comp.Liquids {
		buffer.WriteString(fmt.Sprintf("('%s', '%s', %d)", liquid.Name, liquid.Unit, liquid.Volume))

		// Если это последний элемент массива, то "," в конце не ставим
		if liquid == comp.Liquids[len(comp.Liquids)-1] {
			break
		}

		buffer.WriteString(",")
	}

	buffer.WriteString("] AS Liquid []), CAST(ARRAY[")

	for _, solidBulk := range comp.SolidsBulk {
		buffer.WriteString(fmt.Sprintf("('%s', '%s', %d)", solidBulk.Name, solidBulk.Unit, solidBulk.Volume))

		// Если это последний элемент массива, то "," в конце не ставим
		if solidBulk == comp.SolidsBulk[len(comp.SolidsBulk)-1] {
			break
		}

		buffer.WriteString(",")
	}

	buffer.WriteString("] AS Solid_bulk []), CAST(ARRAY[")

	for _, solidUnit := range comp.SolidsUnit {
		buffer.WriteString(fmt.Sprintf("('%s', %d)", solidUnit.Name, solidUnit.Volume))

		// Если это последний элемент массива, то "," в конце не ставим
		if solidUnit == comp.SolidsUnit[len(comp.SolidsUnit)-1] {
			break
		}

		buffer.WriteString(",")
	}

	buffer.WriteString("] AS Solid_unit [])) AS Composition)")

	return buffer.String()
}

func DecodeDrinkComposition(comp string) (menu.Composition, error) {
	var composition menu.Composition

	comp = strings.ReplaceAll(comp, "\"", "")
	comp = strings.ReplaceAll(comp, "\\", "")

	err := UnmarshalQueryRow(comp, &composition)
	if err != nil {
		return composition, err
	}

	return composition, nil
}

// Not used
func DecodeMenuRequest(dr string) (map[string][]menu.Drink, error) {
	drinks := make(map[string][]menu.Drink, 0)

	dr = strings.ReplaceAll(dr, "\"", "")
	dr = strings.ReplaceAll(dr, "\\", "")

	err := UnmarshalQueryRow(dr, drinks)
	if err != nil {
		return drinks, err
	}

	return drinks, nil
}
