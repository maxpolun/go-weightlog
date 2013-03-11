-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE exersizes (
	name VARCHAR(128) PRIMARY KEY
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE exersizes;
