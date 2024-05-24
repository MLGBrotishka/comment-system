CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    comments_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    name TEXT NOT NULL CHECK (length(name) <= 255),
    text TEXT NOT NULL CHECK (length(text) <= 10000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);