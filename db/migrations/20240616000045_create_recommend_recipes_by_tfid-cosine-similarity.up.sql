-- Compute TF-IDF for ingredients, description words, and new features
CREATE OR REPLACE FUNCTION compute_recipe_tfidf()
RETURNS TABLE (
    recipe_id UUID,
    term TEXT,
    tfidf FLOAT
) AS $$
BEGIN
RETURN QUERY
    WITH recipe_terms AS (
        SELECT
            r.id AS recipe_id,
            unnest(string_to_array(lower(r.description), ' ')) AS term
        FROM
            recipes r
        UNION ALL
        SELECT
            r.id AS recipe_id,
            lower(i.name) AS term
        FROM
            recipes r
            JOIN recipe_ingredients ri ON r.id = ri.recipe_id
            JOIN ingredients i ON ri.ingredient_id = i.id
        UNION ALL
        SELECT
            r.id AS recipe_id,
            lower(rd.meal_type) AS term
        FROM
            recipes r
            JOIN recipe_details rd ON r.id = rd.recipe_id
        UNION ALL
        SELECT
            r.id AS recipe_id,
            lower(rd.food_class) AS term
        FROM
            recipes r
            JOIN recipe_details rd ON r.id = rd.recipe_id
        UNION ALL
        SELECT
            r.id AS recipe_id,
            lower(rd.region) AS term
        FROM
            recipes r
            JOIN recipe_details rd ON r.id = rd.recipe_id
        UNION ALL
        SELECT
            r.id AS recipe_id,
            lower(rd.spiciness_level) AS term
        FROM
            recipes r
            JOIN recipe_details rd ON r.id = rd.recipe_id
        UNION ALL
        SELECT
            r.id AS recipe_id,
            unnest(string_to_array(lower(rd.main_ingredients), ', ')) AS term
        FROM
            recipes r
            JOIN recipe_details rd ON r.id = rd.recipe_id
        UNION ALL
        SELECT
            r.id AS recipe_id,
            lower(rd.cooking_method) AS term
        FROM
            recipes r
            JOIN recipe_details rd ON r.id = rd.recipe_id
    ),
    tf AS (
        SELECT
            rt.recipe_id,
            rt.term,
            COUNT(*)::FLOAT / (SELECT COUNT(*) FROM recipe_terms rt2 WHERE rt2.recipe_id = rt.recipe_id) AS tf
        FROM
            recipe_terms rt
        GROUP BY
            rt.recipe_id, rt.term
    ),
    idf AS (
        SELECT
            rt.term,
            LOG((SELECT COUNT(DISTINCT rt.recipe_id) FROM recipe_terms rt) / COUNT(DISTINCT rt.recipe_id)) AS idf
        FROM
            recipe_terms rt
        GROUP BY
            rt.term
    )
SELECT
    tf.recipe_id,
    tf.term,
    tf.tf * idf.idf AS tfidf
FROM
    tf
        JOIN idf ON tf.term = idf.term;
END;
$$ LANGUAGE plpgsql;

-- Create a materialized view to store combined TF-IDF vectors
CREATE MATERIALIZED VIEW IF NOT EXISTS recipe_vectors AS
SELECT
    recipe_id,
    array_agg(tfidf ORDER BY term) AS tfidf_vector
FROM
    compute_recipe_tfidf()
GROUP BY
    recipe_id;

-- Index the materialized view for faster lookups
CREATE INDEX IF NOT EXISTS idx_recipe_vectors ON recipe_vectors USING GIN (tfidf_vector);

-- Compute cosine similarity between recipes
CREATE OR REPLACE FUNCTION compute_cosine_similarity(recipe_id1 UUID, recipe_id2 UUID)
RETURNS FLOAT AS $$
DECLARE
vec1 FLOAT[];
    vec2 FLOAT[];
    dot_product FLOAT := 0;
    norm1 FLOAT := 0;
    norm2 FLOAT := 0;
BEGIN
SELECT tfidf_vector INTO vec1 FROM recipe_vectors WHERE recipe_id = recipe_id1;
SELECT tfidf_vector INTO vec2 FROM recipe_vectors WHERE recipe_id = recipe_id2;

FOR i IN 1..array_length(vec1, 1) LOOP
        dot_product := dot_product + (vec1[i] * vec2[i]);
        norm1 := norm1 + (vec1[i] * vec1[i]);
        norm2 := norm2 + (vec2[i] * vec2[i]);
END LOOP;

    IF norm1 = 0 OR norm2 = 0 THEN
        RETURN 0;
ELSE
        RETURN dot_product / (sqrt(norm1) * sqrt(norm2));
END IF;
END;
$$ LANGUAGE plpgsql;

-- Recommend recipes based on cosine similarity and meal type
CREATE OR REPLACE FUNCTION recommend_recipes(p_user_id UUID, p_limit INT, p_meal_type TEXT)
RETURNS TABLE (
    recipe_id UUID,
    similarity FLOAT
) AS $$
BEGIN
RETURN QUERY
    WITH user_likes AS (
        SELECT
            unnest(recipe_ids) AS recipe_id
        FROM
            user_preferences
        WHERE
            user_id = p_user_id
    ),
    candidate_recipes AS (
        SELECT
            rv2.recipe_id AS recipe_id,
            MAX(compute_cosine_similarity(ul.recipe_id, rv2.recipe_id)) AS similarity
        FROM
            user_likes ul
            JOIN recipe_vectors rv1 ON ul.recipe_id = rv1.recipe_id
            JOIN recipe_vectors rv2 ON rv1.recipe_id <> rv2.recipe_id
        GROUP BY
            rv2.recipe_id
    ),
    filtered_recipes AS (
        SELECT
            cr.recipe_id,
            cr.similarity
        FROM
            candidate_recipes cr
            JOIN recipe_details rd ON cr.recipe_id = rd.recipe_id
        WHERE
            rd.meal_type = p_meal_type
    )
SELECT
    fr.recipe_id,
    fr.similarity
FROM
    filtered_recipes fr
WHERE
    fr.similarity IS NOT NULL
ORDER BY
    fr.similarity DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;
