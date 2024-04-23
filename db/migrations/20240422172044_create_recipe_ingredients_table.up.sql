CREATE TABLE recipe_ingredients (
                                    recipe_id UUID REFERENCES recipes(id),
                                    ingredient_id UUID REFERENCES ingredients(id),
                                    quantity TEXT,
                                    PRIMARY KEY (recipe_id, ingredient_id)
);
