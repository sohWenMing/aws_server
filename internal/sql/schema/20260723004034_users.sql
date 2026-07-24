-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT UNIQUE NOT NULL DEFAULT 'unset',
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
