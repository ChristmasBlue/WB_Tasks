-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments(
                                       id SERIAL PRIMARY KEY,
                                       parent_id INT REFERENCES comments(id) ON DELETE CASCADE,
    text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    search_vector TSVECTOR
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_comments_search ON comments USING GIN(search_vector);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION comments_search_vector_update() RETURNS trigger AS $func$
BEGIN
  NEW.search_vector := to_tsvector('russian', coalesce(NEW.text, ''));
RETURN NEW;
END;
$func$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER tsvectorupdate BEFORE INSERT
    ON comments FOR EACH ROW EXECUTE FUNCTION comments_search_vector_update();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd