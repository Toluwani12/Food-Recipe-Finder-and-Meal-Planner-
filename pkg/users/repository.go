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

func (r Repository) save(ctx context.Context, data AddRequest) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", data.Email)

	if user != (User{}) {
		return nil, liberror.ErrEmailExists
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "db.GetContext failed")
	}

	rows, err := r.db.NamedQueryContext(ctx, `INSERT INTO users ( username, email, password_hash) VALUES (:username, :email, :password_hash) RETURNING *`, data)
	if err != nil {
		return nil, errors.Wrap(err, "Db.NamedQueryContext")
	}

	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "Rows.StructScan")
		}
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
