
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO exersizes (name) VALUES 
	('squat'), 
	('benchpress'),
	('overheadpress'),
	('deadlift'),
	('barbellrow'),
	('powerclean'),
	('snatch'),
	('chinup'),
	('pullup'),
	('curl')
;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
TRUNCATE exersizes CASCADE;
