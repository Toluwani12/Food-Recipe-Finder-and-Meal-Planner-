package user_preference

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

func (r Repository) add(ctx context.Context, id string, req AddRequest) error {

	_, err := r.db.ExecContext(ctx, "INSERT INTO user_preferences (user_id, vegetarian, gluten_free, cuisine_preference, disliked_ingredients, additional_preferences, dietary_goals) VALUES ($1, $2, $3, $4, $5, $6, $7)", id, req.Vegetarian, req.GlutenFree, req.CuisinePreference, req.DislikedIngredients, req.AdditionalPreferences, req.DietaryGoals)
	if err != nil {
		return errors.Wrap(err, "db.ExecContext failed")
	}
	return nil
}

func (r Repository) delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM user_preferences WHERE user_id = $1", id)
	return errors.Wrap(err, "db.ExecContext failed")
}

func (r Repository) update(ctx context.Context, id string, data AddRequest) error {
	row, err := r.db.ExecContext(ctx, "UPDATE user_preferences SET vegetarian = $1, gluten_free = $2, cuisine_preference = $3, disliked_ingredients = $4, additional_preferences = $5 WHERE user_id = $6", data.Vegetarian, data.GlutenFree, data.CuisinePreference, data.DislikedIngredients, data.AdditionalPreferences, id)
	if err != nil {
		return errors.Wrap(err, "db.ExecContext failed")
	}
	if count, _ := row.RowsAffected(); count == 0 {
		return liberror.ErrNotFound
	}

	return errors.Wrap(err, "db.ExecContext failed")
}

func (r Repository) get(ctx context.Context, id string) (*UserPreference, error) {
	var userPreference UserPreference
	err := r.db.GetContext(ctx, &userPreference, "SELECT * FROM user_preferences WHERE user_id = $1", id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, liberror.ErrNotFound
	}
	return &userPreference, errors.Wrap(err, "db.GetContext failed")
}
