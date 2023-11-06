package user_db

import (
	"context"
	"errors"
	"fmt"
	"restapi/internal/domain/user"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/hash"
	"restapi/pkg/logging"
	repeatable "restapi/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	hasher hash.PasswordHasher
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, dto user.CreateUserDTO) (usr user.User, err error) {
	q := `	
	INSERT INTO users 
		(name, login, password_hash, one_time_code)
	VALUES 
		($1, $2, $3, $4) 
	RETURNING id, name, login, password_hash
	` // CHECK REQUEST!!! //
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	passwordHash, err := r.hasher.Hash(dto.Password)
	if err != nil {
		return user.User{}, err
	}

	if err := r.client.QueryRow(ctx, q, dto.Name, dto.Login, passwordHash, "").Scan(&usr.ID,
		&usr.Name, &usr.Login, &usr.PasswordHash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return user.User{}, newErr
		}

		return user.User{}, err
	}

	return usr, nil
}

// Удаляет юзера и все связанные с ним записи в sessions
func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
	DELETE FROM 
		users 
	WHERE 
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, id)
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
		err := fmt.Errorf("database deleting error: user not found")
		return err
	}

	return nil
}

func (r *repository) GetByCredentials(ctx context.Context, login, passwordHash string) (user.User, error) {
	q := `
	SELECT 
		id, name, login, password_hash, one_time_code
	FROM 
		public.users
	WHERE 
		(login = $1) AND (password_hash = $2)
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var usr user.User

	err := r.client.QueryRow(ctx, q, login, passwordHash).Scan(&usr.ID, &usr.Name, &usr.Login, &usr.PasswordHash,
		&usr.OneTimeCode)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return user.User{}, newErr
		}

		return user.User{}, err
	}

	return usr, nil
}

func (r *repository) GetByUUID(ctx context.Context, userID string) (user.User, error) {
	q := `
	SELECT 
		name, login, password_hash, one_time_code
	FROM 
		public.users
	WHERE 
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var usr user.User

	err := r.client.QueryRow(ctx, q, userID).Scan(&usr.Name, &usr.Login, &usr.PasswordHash, &usr.OneTimeCode)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return user.User{}, newErr
		}

		return user.User{}, err
	}

	usr.ID = userID

	return usr, nil
}

func (r *repository) FindAll(ctx context.Context) ([]user.User, error) {
	// в следующий раз тут будет множество сортировок и опций

	q := `
	SELECT 
		id, name, login, password_hash, one_time_code
	FROM 
		public.users
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]user.User, 0)

	for rows.Next() {
		var usr user.User

		err = rows.Scan(&usr.ID, &usr.Name, &usr.Login, &usr.PasswordHash, &usr.OneTimeCode)
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

		users = append(users, usr)
	}

	return users, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (user.User, error) {
	q := `
	SELECT 
		id, name, login, password_hash, one_time_code
	FROM 
		public.users
	WHERE 
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var usr user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Name, &usr.Login,
		&usr.PasswordHash, &usr.OneTimeCode)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return user.User{}, newErr
		}

		return user.User{}, err
	}

	return usr, nil
}

func (r *repository) Update(ctx context.Context, user user.User) error {
	q := `
	UPDATE users
	SET 
		name = $2, login = $3, password_hash = $4, one_time_code = $5
	WHERE 
		id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, user.ID, user.Name, user.Login, user.PasswordHash, user.OneTimeCode)
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
		err := fmt.Errorf("database updating error: user not found")
		return err
	}

	return nil
}

func (r *repository) PartUpdate(ctx context.Context, dto user.PartUpdateUserDTO) error {
	var q string

	switch dto.Key {
	case "name":
		{
			q = `
			UPDATE users
			SET 
				name = $2
			WHERE 
				id = $1
		`
		}
	case "login":
		{
			q = `
			UPDATE users
			SET 
				login = $2
			WHERE 
				id = $1
			`
		}
	case "password_hash":
		{
			q = `
			UPDATE users
			SET 
				password_hash = $2
			WHERE 
				id = $1
		`
		}
	case "one_time_code":
		{
			q = `
			UPDATE users
			SET 
				one_time_code = $2
			WHERE 
				id = $1
			`
		}
	}

	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	ct, err := r.client.Exec(ctx, q, dto.ID, dto.Value)
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
		err := fmt.Errorf("database updating error: user not found")
		return err
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger, hasher hash.PasswordHasher) user.Repository {
	return &repository{
		client: client,
		logger: logger,
		hasher: hasher,
	}
}
