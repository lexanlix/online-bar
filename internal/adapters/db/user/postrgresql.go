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
	"strings"

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
	RETURNING true AS is_deleted
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	var isDeleted bool

	err := r.client.QueryRow(ctx, q, id).Scan(&isDeleted)
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
		return user.User{}, err
	}

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
			return nil, err
		}

		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
	RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", repeatable.FormatQuery(q)))

	err := r.client.QueryRow(ctx, q, user.ID, user.Name, user.Login, user.PasswordHash,
		user.OneTimeCode).Scan(&user.ID)
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

func NewRepository(client postgresql.Client, logger *logging.Logger, hasher hash.PasswordHasher) user.Repository {
	return &repository{
		client: client,
		logger: logger,
		hasher: hasher,
	}
}
