CREATE TABLE recipes (
                         id TEXT PRIMARY KEY,
                         name TEXT NOT NULL,
                         cooking_time TEXT,
                         instructions TEXT,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
