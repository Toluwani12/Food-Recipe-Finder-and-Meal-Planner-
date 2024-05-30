package ingredient

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

func (r Repository) get(ctx context.Context, id string) (*Ingredient, error) {
	var ingredient Ingredient

	err := r.db.GetContext(ctx, &ingredient, "SELECT * FROM ingredients WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &ingredient, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) getByName(ctx context.Context, name string) (*Ingredient, error) {
	var ingredient Ingredient

	err := r.db.GetContext(ctx, &ingredient, "SELECT * FROM ingredients WHERE name = $1", name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &ingredient, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) save(ctx context.Context, data Ingredient) (*Ingredient, error) {
	var ingredient *Ingredient
	err := r.db.GetContext(ctx, ingredient, "SELECT * FROM ingredients WHERE name = $1", data.Name)

	if ingredient != nil || !errors.Is(err, sql.ErrNoRows) {
		return nil, liberror.New("name already exist", http.StatusBadRequest)
	}
	if err != nil {
		return nil, errors.Wrap(err, "GetContext")
	}

	res, err := r.db.NamedExecContext(ctx, `INSERT INTO ingredients (name, id, alternative, quantity, created_at, updated_at) VALUES (:name, :id, :alternative, :quantity, :created_at, :updated_at)`, data)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ingredients where id = $1`, id)
	return id, errors.Wrap(err, "ExecContext")
}

func (r Repository) update(ctx context.Context, data Ingredient) (*Ingredient, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE ingredients SET name = $1, quantity = $2,  alternative = $3, updated_at = $4 WHERE id = $5", data.Name, data.Quantity, data.Alternatives, data.UpdatedAt, data.ID)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) list(ctx context.Context) ([]Ingredient, error) {
	var ingredients []Ingredient
	err := r.db.GetContext(ctx, &ingredients, "SELECT * FROM ingredients")

	return ingredients, errors.Wrap(err, "GetContext")
}
