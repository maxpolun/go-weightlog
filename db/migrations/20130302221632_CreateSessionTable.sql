
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE sessions (
	id CHAR(88) PRIMARY KEY,
	user_id serial references users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE sessions;
