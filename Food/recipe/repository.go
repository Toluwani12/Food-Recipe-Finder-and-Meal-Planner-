package recipe

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

func (r Repository) getByReference(ref string) (*Recipe, error) {
	var recipe Recipe

	// Use Get to query and automatically scan the result into the struct
	err := r.db.Get(&recipe, "SELECT * FROM recipes WHERE reference = $1", ref)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &recipe, nil
}

func (r Repository) save(data Recipe) (*Recipe, error) {
	data.Id = uuid.NewString()
	_, err := r.db.NamedExec(`INSERT INTO recipes (name, id) VALUES (:name, :id)`, data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r Repository) delete(data Recipe) (*Recipe, error) {
	data.Name = uuid.NewString()
	_, err := r.db.Exec(`DELETE FROM recipes VALUES (:name)`, data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
