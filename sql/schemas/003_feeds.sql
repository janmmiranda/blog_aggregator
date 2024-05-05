-- +goose Up
CREATE TABLE feeds (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(url)
);

-- +goose Down
DROP TABLE feeds;