-- Create the user_preferences table with a new column for recipe names
CREATE TABLE user_preferences (
                                  user_id UUID PRIMARY KEY REFERENCES users(id),
                                  recipe_ids UUID[]  -- Array of recipe IDs
);

-- Create the trigger function
CREATE OR REPLACE FUNCTION validate_recipe_ids()
RETURNS TRIGGER AS $$
BEGIN
    -- If the recipe_ids array is empty, do nothing
    IF array_length(NEW.recipe_ids, 1) IS NULL THEN
        RETURN NEW;
END IF;

    -- Ensure each recipe_id exists in the recipes table
    PERFORM 1 FROM recipes WHERE id = ANY(NEW.recipe_ids);
    IF NOT FOUND THEN
        RAISE EXCEPTION 'One or more recipe IDs are invalid';
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER validate_recipe_ids_trigger
    BEFORE INSERT OR UPDATE ON user_preferences
                         FOR EACH ROW EXECUTE FUNCTION validate_recipe_ids();

