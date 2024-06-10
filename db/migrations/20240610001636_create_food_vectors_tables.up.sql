-- Create a new table for storing food vectors
CREATE TABLE food_vectors (
                              id SERIAL PRIMARY KEY,
                              food_name TEXT UNIQUE,
                              vector FLOAT8[]
);

-- SQL function to compute the average vector
CREATE OR REPLACE FUNCTION average_vector(food_names TEXT[])
RETURNS FLOAT8[] AS $$
DECLARE
avg_vector FLOAT8[];
    temp_vector FLOAT8[];
    i INT;
    j INT;
BEGIN
SELECT array_fill(0.0, array[vector_length]) INTO avg_vector
FROM (
         SELECT array_length(vector, 1) AS vector_length
         FROM food_vectors
         WHERE food_name = food_names[1]
             LIMIT 1
     ) subquery;

FOR i IN 1..array_length(food_names, 1) LOOP
SELECT vector INTO temp_vector
FROM food_vectors
WHERE food_name = food_names[i];

FOR j IN 1..array_length(temp_vector, 1) LOOP
            avg_vector[j] := avg_vector[j] + temp_vector[j];
END LOOP;
END LOOP;

FOR j IN 1..array_length(avg_vector, 1) LOOP
        avg_vector[j] := avg_vector[j] / array_length(food_names, 1);
END LOOP;

RETURN avg_vector;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION cosine_similarity(vec1 FLOAT8[], vec2 FLOAT8[])
RETURNS FLOAT8 AS $$
DECLARE
dot_product FLOAT8 := 0;
    magnitude1 FLOAT8 := 0;
    magnitude2 FLOAT8 := 0;
    i INT;
BEGIN
FOR i IN 1..array_length(vec1, 1) LOOP
        dot_product := dot_product + (vec1[i] * vec2[i]);
        magnitude1 := magnitude1 + (vec1[i] * vec1[i]);
        magnitude2 := magnitude2 + (vec2[i] * vec2[i]);
END LOOP;
    IF magnitude1 = 0 OR magnitude2 = 0 THEN
        RETURN 0;
END IF;
RETURN dot_product / (SQRT(magnitude1) * SQRT(magnitude2));
END;
$$ LANGUAGE plpgsql;
