-- Create the likes table to store liked recipes
CREATE TABLE likes (
                       user_id UUID REFERENCES users(id),
                       recipe_id UUID REFERENCES recipes(id),
                       PRIMARY KEY (user_id, recipe_id)
);