CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER REFERENCES comments(id),
    level INTEGER NOT NULL CHECK (level >= 1),
    user_id INTEGER,
    post_id INTEGER REFERENCES posts(id),
    text TEXT NOT NULL CHECK (length(text) <= 2000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);