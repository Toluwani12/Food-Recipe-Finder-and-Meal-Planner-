package ingredient

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
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

func (r Repository) save(data Ingredient) (*Ingredient, error) {
	//data.ID = uuid.NewString()
	_, err := r.db.NamedExec(`INSERT INTO ingredients (name, id) VALUES (:name, :id)`, data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r Repository) delete(data Ingredient) (*Ingredient, error) {
	data.Name = uuid.NewString()
	_, err := r.db.Exec(`DELETE FROM ingredients VALUES (:name)`, data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
