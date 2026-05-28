-- +goose Up
CREATE TABLE IF NOT EXISTS urls(
                                   id SERIAL PRIMARY KEY NOT NULL,
                                   short_url TEXT NOT NULL UNIQUE,
                                   url TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS urls;