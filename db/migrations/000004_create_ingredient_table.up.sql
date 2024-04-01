CREATE TABLE ingredients (
                             id TEXT PRIMARY KEY,
                             name TEXT NOT NULL,
                             alternative TEXT,
                             quantity TEXT,
                             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);