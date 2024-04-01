package recipe

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

func (r Repository) get(ctx context.Context, id string) (*Recipe, error) {
	var recipe Recipe

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &recipe, "SELECT * FROM recipes WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &recipe, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) getByName(ctx context.Context, name string) (*Recipe, error) {
	var recipe Recipe

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &recipe, "SELECT * FROM recipes WHERE name = $1", name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &recipe, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) save(ctx context.Context, data Recipe) (*Recipe, error) {
	var recipe *Recipe
	err := r.db.GetContext(ctx, recipe, "SELECT * FROM recipes WHERE name = $1", data.Name)

	if recipe != nil || !errors.Is(err, sql.ErrNoRows) {
		return nil, liberror.New("name already exist", http.StatusBadRequest)
	}
	if err != nil {
		return nil, errors.Wrap(err, "GetContext")
	}

	res, err := r.db.NamedExecContext(ctx, `INSERT INTO recipes (id, name, cooking_time, instructions, created_at, updated_at) VALUES (:id, :name, :cooking_time, :instructions, :created_at, :updated_at)`, data)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM recipes VALUES (:id)`, id)
	return id, errors.Wrap(err, "ExecContext")
}

func (r Repository) list(ctx context.Context) ([]Recipe, error) {
	var recipes []Recipe
	err := r.db.GetContext(ctx, &recipes, "SELECT * FROM recipes")

	return recipes, errors.Wrap(err, "GetContext")
}

func (r Repository) update(ctx context.Context, data Recipe) (*Recipe, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE recipes SET name = $1, cooking_time = $2, instructions = $3, updated_at = $4 WHERE id = $5", data.Name, data.CookingTime, data.Instructions, data.UpdateAt, data.Id)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}
