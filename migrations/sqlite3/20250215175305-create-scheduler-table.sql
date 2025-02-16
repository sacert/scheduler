
-- +migrate Up
CREATE TABLE scheduler (
	id TEXT PRIMARY KEY, -- // uuid 
	data JSON NOT NULL,  -- data json
	schedule_time INTEGER NOT NULL,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL
);


-- +migrate Down
DROP TABLE scheduler;
