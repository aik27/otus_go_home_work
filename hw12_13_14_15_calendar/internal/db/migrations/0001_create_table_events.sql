-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
     id SERIAL PRIMARY KEY,
     title VARCHAR(100) NOT NULL,
     description TEXT NOT NULL,
     start_date TIMESTAMPTZ NOT NULL,
     end_date TIMESTAMPTZ NOT NULL,
     user_id INTEGER NOT NULL,
     remind_in INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd