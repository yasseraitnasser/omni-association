BEGIN;

-- ENUMs
CREATE TYPE member_role AS ENUM (
	'president',
	'vice-president',
	'treasurer',
	'assistant-treasurer',
	'general-secretary',
	'assistant-general-secretary',
	'advisor',
	'member'
);

CREATE TYPE project_status AS ENUM (
	'planning',
	'active',
	'completed',
	'archived'
);

CREATE TYPE project_roles AS ENUM (
	'contributer',
	'project-lead'
);
-- ENUMs

-- tables
CREATE TABLE IF NOT EXISTS members (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL,
	email		VARCHAR(255) UNIQUE NOT NULL,
	role		member_role NOT NULL DEFAULT 'member',
	paid_fee	BOOLEAN DEFAULT FALSE,
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS association (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL DEFAULT 'omni-association',
	current_fee	INTEGER NOT NULL,
	currency	VARCHAR(3) DEFAULT 'MAD',
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS membership_payments (
	id		SERIAL PRIMARY KEY,
	member_id	INTEGER REFERENCES members(id) ON DELETE CASCADE,
	amount_paid	INTEGER NOT NULL,
	year		INTEGER NOT NULL,
	paid_at		TIMESTAMP DEFAULT now() NOT NULL,

	CONSTRAINT unique_member_year UNIQUE (member_id, year)
);

CREATE TABLE IF NOT EXISTS projects (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL,
	budget		INTEGER NOT NULL DEFAULT 0,
	status		project_status NOT NULL DEFAULT 'planning',
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS project_members (
	project_id	INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	member_id	INTEGER REFERENCES members(id) ON DELETE CASCADE,
	role_in_project	project_roles NOT NULL DEFAULT 'contributer',
	join_at		TIMESTAMP DEFAULT now() NOT NULL,

	PRIMARY KEY (project_id, member_id)
);
-- tables

COMMIT;
