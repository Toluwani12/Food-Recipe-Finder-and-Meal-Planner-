package users

import (
	liberror "Food/internal/errors"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) get(ctx context.Context, id string) (*User, error) {
	var user User

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &user, errors.Wrap(err, "db.GetContext failed")
}

func (r *Repository) save(ctx context.Context, data AddRequest) (*User, error) {
	var user User

	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "tx.BeginTxx failed")
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

	// Check if the user already exists
	err = tx.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", data.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "tx.GetContext failed")
	}
	if user != (User{}) {
		return nil, liberror.ErrEmailExists
	}

	// Insert the new user
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id`, data.Username, data.Email, data.PasswordHash).Scan(&user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "tx.QueryRowxContext failed")
	}

	// Create a user_preferences entry
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_preferences (user_id, recipe_ids)
		VALUES ($1, '{}')`, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "tx.ExecContext for user_preferences failed")
	}

	return &user, nil
}
func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users VALUES (:id)`, id)
	return id, errors.Wrap(err, "ExecContext")
}

func (r Repository) list(ctx context.Context) (Users, error) {
	var users Users
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users")

	return users, errors.Wrap(err, "GetContext")
}

func (r Repository) findByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &user, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) update(ctx context.Context, id string, data UpdateRequest) (string, error) {

	var u User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = $1 AND id != $2", data.Email, id)
	if u != (User{}) {
		return "", liberror.ErrEmailExists
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", errors.Wrap(err, "db.GetContext failed")
	}

	res, err := r.db.ExecContext(ctx, `UPDATE users SET username = $1, email = $2 WHERE id = $3`, data.Name, data.Email, id)
	if err != nil {
		return "", errors.Wrap(err, "ExecContext")
	}

	if count, err := res.RowsAffected(); count != 1 {
		return "", errors.Wrap(err, "RowsAffected")
	}

	return id, nil

}
