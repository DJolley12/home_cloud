CREATE TYPE migration_state AS ENUM (
	'DOWN',
	'UP'
);

CREATE TABLE IF NOT EXISTS migrations (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(25) NOT NULL DEFAULT '',
	description VARCHAR(50) NOT NULL DEFAULT '',
	migration_state migration_state DEFAULT 'DOWN'
);
