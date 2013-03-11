
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TYPE weight_unit AS ENUM ('lb', 'kg');

CREATE TABLE sets (
	id SERIAL PRIMARY KEY,
	completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	exersize varchar(128) REFERENCES exersizes(name),
	reps INTEGER NOT NULL,
	user_id SERIAL REFERENCES users(id),
	notes TEXT,
	weight INTEGER NOT NULL,
	unit weight_unit NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE sets;
DROP TYPE weight_unit;
