-- +goose Up
CREATE TABLE IF NOT EXISTS redirect_analytics(
                                                 id SERIAL PRIMARY KEY,
                                                 short_url TEXT NOT NULL REFERENCES urls(short_url),
    request_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_agent TEXT NOT NULL
    );


-- +goose Down
DROP TABLE IF EXISTS redirect_analytics;