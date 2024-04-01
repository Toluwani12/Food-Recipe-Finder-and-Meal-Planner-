CREATE TABLE mealplans (
                           id TEXT PRIMARY KEY,
                           date TIMESTAMP NOT NULL,
                           meal_type TEXT NOT NULL,
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);