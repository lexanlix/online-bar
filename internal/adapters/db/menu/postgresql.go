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

func (r *repository) CreateMenu(ctx context.Context, dto menu.CreateMenuDTO, totalCost uint32) (string, error) {
	drinksString := EncodeInsertValue(&dto.Drinks)

	q := fmt.Sprintf(`
	INSERT INTO menu
		(user_id, name, drinks, total_cost)
	VALUES
    	($1, $2, %s, $3)
	RETURNING
    	id
	`, drinksString)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var menuID string

	row := r.client.QueryRow(ctx, q, dto.UserID, dto.Name, totalCost)
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

func (r *repository) DeleteMenu(ctx context.Context, dto menu.DeleteMenuDTO) error {
	q := `
	DELETE FROM 
		menu
	WHERE
    	id = $1
	RETURNING true AS is_deleted
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var isDeleted bool

	err := r.client.QueryRow(ctx, q, dto.ID).Scan(&isDeleted)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 22P02") {
			err := fmt.Errorf("database error: %v", err)
			return err
		}

		err := fmt.Errorf("database error: rows not found")
		return err
	}

	if !isDeleted {
		err := fmt.Errorf("database deleting error: %v", err)
		return err
	}

	return nil
}

func (r *repository) FindMenu(ctx context.Context, dto menu.FindMenuDTO) (menu.Menu, error) {
	q := `
	SELECT
		id, user_id, name, drinks, total_cost
	FROM 
    	menu
	WHERE 
    	id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var mn menu.Menu
	var drinks string

	err := r.client.QueryRow(ctx, q, dto.ID).Scan(&mn.ID, &mn.UserID, &mn.Name, &drinks, &mn.TotalCost)
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

	drinksMap, err := DecodeMenuRequest(drinks)
	if err != nil {
		return menu.Menu{}, err
	}

	mn.Drinks = drinksMap

	return mn, nil
}

func (r *repository) UpdateMenu(ctx context.Context, dto menu.UpdateMenuDTO, totalCost uint32) (string, error) {
	drinksString := EncodeInsertValue(&dto.Drinks)

	q := fmt.Sprintf(`
	UPDATE menu
	SET
		name = $2, drinks = %s, total_cost = $3
	WHERE
		id = $1
	RETURNING id
	`, drinksString)
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var updatedID string
	err := r.client.QueryRow(ctx, q, dto.ID, dto.Name, totalCost).Scan(&updatedID)
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

// Добавляет напиток в массив определенной категории, или в новую созданную категорию
func (r *repository) AddDrink(ctx context.Context, dto menu.AddDrinkDTO) error {
	var q string

	if dto.IsNewCateg {
		drinkString := EncodeInsertDrinkWithCategory(&dto.Drink)

		q = fmt.Sprintf(`
		UPDATE menu
		SET
			drinks[%d] = %s,
		    total_cost = total_cost + %d
		WHERE
			id = $1
		RETURNING id;
		`, dto.CategID, drinkString, dto.Drink.Price)
	} else {
		drinkString := EncodeInsertDrink(&dto.Drink)

		q = fmt.Sprintf(`
		UPDATE menu
		SET
			drinks[%d].drinks[%d] = %s,
			total_cost = total_cost + %d
		WHERE
			id = $1
		RETURNING id
		`, dto.CategID, dto.Drink.ID, drinkString, dto.Drink.Price)
	}

	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	err := r.client.QueryRow(ctx, q, dto.MenuID).Scan(&dto.MenuID)
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

func NewRepository(client postgresql.Client, logger *logging.Logger) menu.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func EncodeInsertValue(drinkGroups *map[string][]menu.Drink) string {
	var buffer bytes.Buffer
	var categories []string

	for category := range *drinkGroups {
		categories = append(categories, category)
	}

	buffer.WriteString("CAST(ARRAY[")

	for category, drinks := range *drinkGroups {
		buffer.WriteString(fmt.Sprintf("('%s', CAST(ARRAY[", category))

		for _, drink := range drinks {
			buffer.WriteString(fmt.Sprintf("(%d, '%s', '%s', '%s', CAST(", drink.ID, drink.Name, drink.Category,
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
				buffer.WriteString(fmt.Sprintf("('%s', %d)", solidUnit.Name, solidUnit.Amount))

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

			// Если это последний элемент массива, то "," в конце не ставим
			if drink.ID == drinks[len(drinks)-1].ID {
				break
			}

			buffer.WriteString(",")
		}

		buffer.WriteString("] AS Drink []))")

		// Если это последний элемент map'a, то "," в конце не ставим
		if category == categories[len(categories)-1] {
			break
		}

		buffer.WriteString(",")
	}

	buffer.WriteString("] AS DrinksGroup [])")
	return buffer.String()
}

func EncodeInsertDrink(drink *menu.Drink) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("CAST((%d, '%s', '%s', '%s', CAST(", drink.ID, drink.Name, drink.Category,
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
		buffer.WriteString(fmt.Sprintf("('%s', %d)", solidUnit.Name, solidUnit.Amount))

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

	buffer.WriteString("]) AS Drink)")

	return buffer.String()
}

func EncodeInsertDrinkWithCategory(drink *menu.Drink) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("CAST(('%s', CAST(ARRAY[", drink.Category))

	buffer.WriteString(fmt.Sprintf("(%d, '%s', '%s', '%s', CAST(", drink.ID, drink.Name, drink.Category,
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
		buffer.WriteString(fmt.Sprintf("('%s', %d)", solidUnit.Name, solidUnit.Amount))

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

	buffer.WriteString("])] AS Drink[])) AS DrinksGroup)")

	return buffer.String()
}

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
