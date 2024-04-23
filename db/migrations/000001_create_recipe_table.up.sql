CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE recipes (
                         id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                         description TEXT,
                         name TEXT NOT NULL,
                         cooking_time TEXT,
                         instructions TEXT,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
