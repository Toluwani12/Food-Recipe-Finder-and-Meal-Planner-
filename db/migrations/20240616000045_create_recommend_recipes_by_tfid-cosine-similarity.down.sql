DROP FUNCTION IF EXISTS recommend_recipes(user_id UUID, l INT);
DROP FUNCTION IF EXISTS compute_cosine_similarity(recipe_id1 UUID, recipe_id2 UUID);
DROP MATERIALIZED VIEW IF EXISTS recipe_vectors;
DROP FUNCTION IF EXISTS compute_recipe_tfidf();
