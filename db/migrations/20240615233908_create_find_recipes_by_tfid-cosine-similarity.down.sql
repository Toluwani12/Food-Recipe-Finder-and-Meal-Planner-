DROP TABLE IF EXISTS tfidf_vectors;
DROP FUNCTION IF EXISTS cosine_similarity(vec1 REAL[], vec2 REAL[]);
DROP FUNCTION IF EXISTS find_recipes_by_cosine_similarity(user_ingredients TEXT[]);