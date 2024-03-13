package ingredient

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) getByReference(ref string) (*Ingredient, error) {
	var ingredient Ingredient

	// Use Get to query and automatically scan the result into the struct
	err := r.db.Get(&ingredient, "SELECT * FROM ingredients WHERE reference = $1", ref)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &ingredient, nil
}

func (r Repository) save(ctx context.Context, data AddRequest) (*Ingredient, error) {
	//data.ID = uuid.NewString()

	_, err := r.db.NamedExecContext(ctx, `INSERT INTO ingredients (name, id) VALUES (:name, :id)`, data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r Repository) delete(ctx context.Context, id string) (*Ingredient, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ingredients VALUES (:id)`, id)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// for this update, it'd need to collect the new data and probably bind it
func (r Repository) update(id string, data AddRequest) (*Ingredient, error) {
	// Construct the update query with newData fields and the id
	_, err := r.db.Exec("UPDATE ingredients SET name = $1, quantity = $2 WHERE id = $3", data.Name, data.Quantity, id)

	if err != nil {
		return err
	}

	return nil
}
