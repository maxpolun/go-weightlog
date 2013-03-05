-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE exersizes (
	id SERIAL PRIMARY KEY,
	name VARCHAR(150)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE exersizes;
