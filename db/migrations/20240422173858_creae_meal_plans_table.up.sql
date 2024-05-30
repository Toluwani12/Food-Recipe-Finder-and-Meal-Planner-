-- Create ENUM type meal_type if it does not exist
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'meal_type') THEN
CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'dinner');
END IF;
END $$;

-- Create ENUM type day_of_week if it does not exist
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'day_of_week') THEN
CREATE TYPE day_of_week AS ENUM ('sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday');
END IF;
END $$;

-- Create the meal_plans table if it does not exist
CREATE TABLE IF NOT EXISTS meal_plans (
                                          id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    day_of_week day_of_week,
    meal_type meal_type,
    recipe_id UUID REFERENCES recipes(id),
    week_start_date DATE NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
