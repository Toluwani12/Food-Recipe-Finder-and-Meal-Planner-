CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE ingredients (
                             id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                             name TEXT NOT NULL,
                             created_at TIMESTAMP,
                             updated_at TIMESTAMP
);