
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	email VARCHAR(80) UNIQUE NOT NULL,
	pw_hash CHAR(60) NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users;
