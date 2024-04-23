CREATE TABLE ingredient_alternatives (
                                         ingredient_id UUID REFERENCES ingredients(id),
                                         alternative_id UUID REFERENCES ingredients(id),
                                         PRIMARY KEY (ingredient_id, alternative_id)
);
