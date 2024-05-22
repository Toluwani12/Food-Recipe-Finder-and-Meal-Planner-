CREATE TABLE recipes (
                         id UUID PRIMARY KEY,
                         description TEXT,
                         name TEXT NOT NULL,
                         cooking_time TEXT,
                         instructions TEXT,
                         img_url TEXT,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE recipes ADD CONSTRAINT unique_recipe_name UNIQUE (name);

