package user_preference

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) get(ctx context.Context, userID string) ([]string, error) {
	var recipeNames []string
	query := `
		SELECT r.name
		FROM recipes r
		JOIN unnest(
			(SELECT recipe_ids
			 FROM user_preferences
			 WHERE user_id = $1)
		) as recipe_id ON r.id = recipe_id`
	err := r.db.SelectContext(ctx, &recipeNames, query, userID)
	if err != nil {
		return nil, err
	}

	return recipeNames, nil
}

func (r *Repository) add(ctx context.Context, tx *sqlx.Tx, userID string, recipeIDs []uuid.UUID) error {
	query := `
        UPDATE user_preferences
        SET recipe_ids = (
            SELECT array_agg(DISTINCT id)
            FROM unnest(array_cat(recipe_ids, $1::uuid[])) AS id
        )
        WHERE user_id = $2`
	_, err := tx.ExecContext(ctx, query, pq.Array(recipeIDs), userID)
	if err != nil {
		return errors.Wrap(err, "exec failed")
	}
	return nil
}

func (r *Repository) remove(ctx context.Context, tx *sqlx.Tx, userID string, recipeIDs []uuid.UUID) error {
	query := `
        UPDATE user_preferences
        SET recipe_ids = (
            SELECT array_agg(id)
            FROM (
                SELECT unnest(recipe_ids) AS id
                EXCEPT
                SELECT unnest($1::uuid[]) AS id_to_remove
            ) AS filtered_ids
        )
        WHERE user_id = $2`
	_, err := tx.ExecContext(ctx, query, pq.Array(recipeIDs), userID)
	if err != nil {
		return errors.Wrap(err, "exec failed")
	}
	return nil
}

func (r *Repository) setLikeStatus(ctx context.Context, userID string, recipeIDs []string, like bool) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
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

	recipeUUIDs := make([]uuid.UUID, len(recipeIDs))
	for i, id := range recipeIDs {
		recipeUUIDs[i] = uuid.MustParse(id)
	}

	if like {
		//Like the recipe
		_, err = tx.ExecContext(ctx, `
			INSERT INTO likes (user_id, recipe_id)
			SELECT $1, unnest($2::uuid[])
			ON CONFLICT DO NOTHING`, userID, pq.Array(recipeUUIDs))
		if err != nil {
			return err
		}

		// Add recipe ID to user preferences
		err = r.add(ctx, tx, userID, recipeUUIDs)
		if err != nil {
			return err
		}
	} else {
		// Unlike the recipe
		_, err = tx.ExecContext(ctx, `
			DELETE FROM likes WHERE user_id = $1 AND recipe_id = ANY($2::uuid[])`, userID, pq.Array(recipeUUIDs))
		if err != nil {
			return err
		}

		// Remove recipe ID from user preferences
		err = r.remove(ctx, tx, userID, recipeUUIDs)
		if err != nil {
			return err
		}
	}

	return err
}
