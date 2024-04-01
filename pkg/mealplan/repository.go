package mealplan

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

func (r Repository) get(ctx context.Context, id string) (*MealPlan, error) {
	var mealPlan MealPlan

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &mealPlan, "SELECT * FROM mealPlans WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &mealPlan, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) getById(ctx context.Context, id string) (*MealPlan, error) {
	var mealPlan MealPlan

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &mealPlan, "SELECT * FROM mealPlans WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &mealPlan, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) save(ctx context.Context, data MealPlan) (*MealPlan, error) {
	var mealPlan *MealPlan
	err := r.db.GetContext(ctx, mealPlan, "SELECT * FROM mealPlans WHERE id = $1", data.Id)

	if mealPlan != nil || !errors.Is(err, sql.ErrNoRows) {
		return nil, liberror.New("name already exist", http.StatusBadRequest)
	}
	if err != nil {
		return nil, errors.Wrap(err, "GetContext")
	}

	res, err := r.db.NamedExecContext(ctx, `INSERT INTO mealPlans (id, date, mealtype) VALUES (:id, :date, :meal_type)`, data)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM mealPlans where id = $1`, id)
	return id, errors.Wrap(err, "ExecContext")
}

// for this update, it'd need to collect the new data and probably bind it
func (r Repository) update(ctx context.Context, data MealPlan) (*MealPlan, error) {
	// Construct the update query with newData fields and the id
	res, err := r.db.ExecContext(ctx, "UPDATE mealPlans SET date = $1,  meal_type = $2 WHERE id = $3", data.Date, data.MealType, data.Id)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) list(ctx context.Context) ([]MealPlan, error) {
	var mealPlans []MealPlan
	err := r.db.GetContext(ctx, &mealPlans, "SELECT * FROM mealPlans")

	return mealPlans, errors.Wrap(err, "GetContext")
}
