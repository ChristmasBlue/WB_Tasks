-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items(
                                    id SERIAL PRIMARY KEY,
                                    name VARCHAR(255) NOT NULL,
    count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
                                    id SERIAL PRIMARY KEY,
                                    name VARCHAR(255) NOT NULL,
    role VARCHAR(10) CHECK (role IN ('manager', 'admin', 'viewer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE items_history (
                               id SERIAL PRIMARY KEY,
                               item_id INT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
                               changed_by_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               changed_column VARCHAR(255) NOT NULL CHECK (changed_column IN ('name', 'count')),
                               changed_from VARCHAR(255) NOT NULL,
                               change_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_items_update()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.name IS DISTINCT FROM NEW.name THEN
    INSERT INTO items_history (item_id, changed_by_id, changed_column, changed_from, change_time)
    VALUES (OLD.id, current_setting('app.current_user_id')::INT, 'name', OLD.name, NOW());
END IF;

  IF OLD.count IS DISTINCT FROM NEW.count THEN
    INSERT INTO items_history (item_id, changed_by_id, changed_column, changed_from, change_time)
    VALUES (OLD.id, current_setting('app.current_user_id')::INT, 'count', OLD.count::VARCHAR, NOW());
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_log_items_update
    AFTER UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_items_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS items_history
-- +goose StatementEnd