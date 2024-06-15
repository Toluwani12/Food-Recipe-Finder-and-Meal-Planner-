DROP FUNCTION IF EXISTS find_recipes_by_jaccard_similarity(TEXT[]);
DROP FUNCTION IF EXISTS find_recipes_by_ingredient_difference(TEXT[]);
DROP MATERIALIZED VIEW IF EXISTS recipe_ingredient_vector;