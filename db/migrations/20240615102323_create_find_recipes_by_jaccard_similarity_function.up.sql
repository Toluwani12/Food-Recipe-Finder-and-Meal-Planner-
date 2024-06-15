-- Create materialized view for recipe ingredient vectors
CREATE MATERIALIZED VIEW recipe_ingredient_vector AS
SELECT
    r.id AS recipe_id,
    r.name AS recipe_name,
    array_agg(i.name ORDER BY i.name) AS ingredients_vector
FROM recipes r
         JOIN recipe_ingredients ri ON r.id = ri.recipe_id
         JOIN ingredients i ON ri.ingredient_id = i.id
GROUP BY r.id;

-- Index the materialized view
CREATE INDEX idx_recipe_ingredient_vector ON recipe_ingredient_vector USING GIN (ingredients_vector);

-- Create function to find recipes by Jaccard similarity
CREATE OR REPLACE FUNCTION find_recipes_by_jaccard_similarity(user_ingredients TEXT[])
    RETURNS TABLE (
                      recipe_id UUID,
                      similarity_score FLOAT
                  ) AS $$
BEGIN
    RETURN QUERY
        WITH user_ingredient_vector AS (
            SELECT user_ingredients AS ingredients_vector
        ),
             ingredient_similarity AS (
                 SELECT
                     r.id AS recipe_id,
                     array_length(ARRAY(
                                          SELECT unnest(riv.ingredients_vector)
                                          INTERSECT
                                          SELECT unnest(u.ingredients_vector)
                                  ), 1) AS intersection_size,
                     array_length(ARRAY(
                                          SELECT unnest(riv.ingredients_vector)
                                          UNION
                                          SELECT unnest(u.ingredients_vector)
                                  ), 1) AS union_size
                 FROM recipes r
                          JOIN recipe_ingredient_vector riv ON r.id = riv.recipe_id, user_ingredient_vector u
             ),
             jaccard_similarity AS (
                 SELECT
                     ingredient_similarity.recipe_id,
                     (ingredient_similarity.intersection_size::float) AS similarity_score
                 FROM ingredient_similarity
                    WHERE ingredient_similarity.union_size > 0 AND ingredient_similarity.intersection_size > 0
             )
        SELECT
            jaccard_similarity.recipe_id,
            jaccard_similarity.similarity_score
        FROM jaccard_similarity
        ORDER BY jaccard_similarity.similarity_score DESC
    LIMIT 10;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION find_recipes_by_ingredient_difference(user_ingredients TEXT[])
    RETURNS TABLE (
                      recipe_id UUID,
                      difference INT
                  ) AS $$
BEGIN
    RETURN QUERY
        WITH recipe_intersections AS (
            SELECT
                r.id AS recipe_id,
                COUNT(DISTINCT i.name) AS intersection_count,
                (SELECT COUNT(*) FROM unnest(user_ingredients)) AS user_ingredient_count,
                (SELECT COUNT(*) FROM recipe_ingredients ri2 WHERE ri2.recipe_id = r.id) AS recipe_ingredient_count
            FROM
                recipes r
                    JOIN recipe_ingredients ri ON r.id = ri.recipe_id
                    JOIN ingredients i ON ri.ingredient_id = i.id
            WHERE
                i.name = ANY(user_ingredients)
            GROUP BY
                r.id
        )
        SELECT
            ri.recipe_id,
            ((ri.recipe_ingredient_count - ri.intersection_count))::integer AS difference
        FROM
            recipe_intersections ri
        WHERE
            ri.intersection_count > 0  -- Ensure at least one matching ingredient
        ORDER BY
            difference ASC;
END;
$$ LANGUAGE plpgsql;

-- Refresh materialized view to keep it up-to-date
REFRESH MATERIALIZED VIEW recipe_ingredient_vector;