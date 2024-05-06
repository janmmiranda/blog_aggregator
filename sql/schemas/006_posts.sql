-- +goose Up
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT,
    url TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id TEXT NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE(url)
);

-- +goose Down
DROP TABLE posts;