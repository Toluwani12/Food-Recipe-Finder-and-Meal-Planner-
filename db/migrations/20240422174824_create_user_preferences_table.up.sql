CREATE TABLE user_preferences (
                                  user_id UUID PRIMARY KEY REFERENCES users(id),
                                  vegetarian BOOLEAN DEFAULT FALSE,
                                  gluten_free BOOLEAN DEFAULT FALSE,
                                  cuisine_preference VARCHAR(255),  -- e.g., "Italian", "Mexican"
                                  disliked_ingredients TEXT[],  -- Array of ingredient names
                                  additional_preferences JSONB,  -- Using JSONB for better performance in PostgreSQL
                                  dietary_goals TEXT
);