package users

import (
	liberror "Food/internal/errors"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) get(ctx context.Context, id string) (*User, error) {
	var user User

	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No user found with the specified ID", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "GetContext: failed to get user by ID")
	}

	return &user, nil
}

func (r *Repository) save(ctx context.Context, data AddRequest) (*User, error) {
	var user User

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "BeginTxx: failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = tx.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", data.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "GetContext: failed to get user by email")
	}
	if user != (User{}) {
		return nil, liberror.New("Email already exists", http.StatusConflict)
	}

	err = tx.QueryRowxContext(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id`, data.Username, data.Email, data.PasswordHash).Scan(&user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "QueryRowxContext: failed to insert new user")
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_preferences (user_id, recipe_ids)
		VALUES ($1, '{}')`, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "ExecContext: failed to insert user preferences")
	}

	return &user, nil
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return id, errors.Wrap(err, "ExecContext: failed to delete user")
	}
	return id, nil
}

func (r Repository) list(ctx context.Context) (Users, error) {
	var users Users
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, errors.Wrap(err, "SelectContext: failed to list users")
	}
	return users, nil
}

func (r Repository) findByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No user found with the specified email", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "GetContext: failed to get user by email")
	}

	return &user, nil
}

func (r Repository) update(ctx context.Context, id string, data UpdateRequest) (string, error) {
	var u User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = $1 AND id != $2", data.Email, id)
	if u != (User{}) {
		return "", liberror.New("Email already exists", http.StatusConflict)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", errors.Wrap(err, "GetContext: failed to get user by email and ID")
	}

	res, err := r.db.ExecContext(ctx, `UPDATE users SET username = $1, email = $2 WHERE id = $3`, data.Name, data.Email, id)
	if err != nil {
		return "", errors.Wrap(err, "ExecContext: failed to update user")
	}

	if count, err := res.RowsAffected(); err != nil {
		return "", errors.Wrap(err, "RowsAffected: failed to get rows affected count")
	} else if count != 1 {
		return "", errors.New("No rows affected")
	}

	return id, nil
}
