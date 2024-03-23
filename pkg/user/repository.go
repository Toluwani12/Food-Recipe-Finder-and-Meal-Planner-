package user

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

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &user, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) save(ctx context.Context, data User) (*User, error) {
	var user *User
	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE name = $1", data.Email)

	if user != nil || !errors.Is(err, sql.ErrNoRows) {
		return nil, liberror.New("email already exist", http.StatusBadRequest)
	}
	if err != nil {
		return nil, errors.Wrap(err, "GetContext")
	}
	res, err := r.db.NamedExecContext(ctx, `INSERT INTO users (id, email, password) VALUES (:id, :email, :password)`, data)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users VALUES (:id)`, id)
	return id, errors.Wrap(err, "ExecContext")
}

func (r Repository) list(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.GetContext(ctx, &users, "SELECT * FROM users")

	return users, errors.Wrap(err, "GetContext")
}
