CREATE TABLE tfidf_vectors (
                               recipe_id UUID PRIMARY KEY,
                               ingredient_vector REAL[]
);

-- Populate tfidf_vectors table
WITH ingredient_count AS (
    SELECT
        ri.recipe_id,
        i.name,
        COUNT(*) AS ingredient_count
    FROM
        recipe_ingredients ri
            JOIN ingredients i ON ri.ingredient_id = i.id
    GROUP BY
        ri.recipe_id, i.name
),
     ingredient_total AS (
         SELECT
             recipe_id,
             SUM(ingredient_count) AS total_count
         FROM
             ingredient_count
         GROUP BY
             recipe_id
     ),
     ingredient_tfidf AS (
         SELECT
             ic.recipe_id,
             ic.name,
             ic.ingredient_count::REAL / it.total_count AS tfidf
         FROM
             ingredient_count ic
                 JOIN ingredient_total it ON ic.recipe_id = it.recipe_id
     )
INSERT INTO tfidf_vectors (recipe_id, ingredient_vector)
SELECT
    itf.recipe_id,
    ARRAY_AGG(itf.tfidf ORDER BY itf.name)
FROM
    ingredient_tfidf itf
GROUP BY
    itf.recipe_id;


CREATE OR REPLACE FUNCTION cosine_similarity(vec1 REAL[], vec2 REAL[])
RETURNS REAL AS $$
BEGIN
RETURN (
           SELECT SUM((v1.val * v2.val))
           FROM unnest(vec1) WITH ORDINALITY AS v1(val, idx)
                    JOIN unnest(vec2) WITH ORDINALITY AS v2(val, idx)
                         ON v1.idx = v2.idx
       ) / (
           (SELECT SQRT(SUM(v1.val ^ 2)) FROM unnest(vec1) WITH ORDINALITY AS v1(val, idx))
               * (SELECT SQRT(SUM(v2.val ^ 2)) FROM unnest(vec2) WITH ORDINALITY AS v2(val, idx))
           );
END;
$$ LANGUAGE plpgsql;

   CREATE OR REPLACE FUNCTION find_recipes_by_cosine_similarity(user_ingredients TEXT[])
RETURNS TABLE (
    recipe_id UUID,
    similarity_score REAL
) AS $$
DECLARE
user_vector REAL[];
BEGIN
    -- Compute the TF-IDF vector for the user ingredients
    user_vector := (
        SELECT ARRAY_AGG(tfidf ORDER BY name)
        FROM (
            SELECT
                i.name,
                COUNT(*)::REAL / (SELECT COUNT(*) FROM unnest(user_ingredients)) AS tfidf
            FROM
                unnest(user_ingredients) AS i(name)
            GROUP BY
                i.name
        ) AS user_tfidf
    );

RETURN QUERY
SELECT
    tf.recipe_id,
    cosine_similarity(tf.ingredient_vector, user_vector) AS similarity_score
FROM
    tfidf_vectors tf
ORDER BY
    similarity_score DESC;
END;
$$ LANGUAGE plpgsql;
