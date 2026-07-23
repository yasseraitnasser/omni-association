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
	'contributor',
	'project-lead'
);

CREATE TYPE funding_source AS ENUM (
	'member',
	'donor_individual',
	'government',
	'external_association'
);

CREATE TYPE transaction_type AS ENUM ('income', 'expense');
CREATE TYPE payment_method AS ENUM ('bank_transfer', 'check', 'cash');

-- tables
CREATE TABLE IF NOT EXISTS association (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL DEFAULT 'omni-association',
	current_fee	INTEGER NOT NULL,
	currency	VARCHAR(3) DEFAULT 'MAD',
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS members (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL,
	email		VARCHAR(255) UNIQUE NOT NULL,
	password	TEXT DEFAULT NULL,
	invite_token	VARCHAR(255) UNIQUE DEFAULT NULL,
	invite_expiry	TIMESTAMP DEFAULT NULL,
	role		member_role NOT NULL DEFAULT 'member',
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS membership_payments (
	id		SERIAL PRIMARY KEY,
	member_id	INTEGER REFERENCES members(id) ON DELETE CASCADE,
	amount_paid	INTEGER NOT NULL,
	year		INTEGER NOT NULL,
	receipt_url	TEXT NOT NULL,
	paid_at		TIMESTAMP DEFAULT now() NOT NULL,

	CONSTRAINT unique_member_year UNIQUE (member_id, year)
);

CREATE TABLE IF NOT EXISTS projects (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL,
	description	TEXT NOT NULL,
	budget		INTEGER NOT NULL DEFAULT 0,
	status		project_status NOT NULL DEFAULT 'planning',
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS project_committees (
	project_id	INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	member_id	INTEGER REFERENCES members(id) ON DELETE CASCADE,
	role_in_project	project_roles NOT NULL DEFAULT 'contributor',
	join_at		TIMESTAMP DEFAULT now() NOT NULL,

	PRIMARY KEY (project_id, member_id)
);

CREATE TABLE IF NOT EXISTS transactions (
	id			SERIAL PRIMARY KEY,
	project_id		INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	recorded_by		INTEGER REFERENCES members(id) ON DELETE SET NULL,
	type			transaction_type NOT NULL,
	source			funding_source NOT NULL,
	donor_member_id		INTEGER REFERENCES members(id) ON DELETE SET NULL,
	external_entity_name	TEXT DEFAULT NULL,
	amount			INTEGER NOT NULL,
	payment_method		payment_method NOT NULL,

	-- document tracking
	proof_doc_url		TEXT NOT NULL,
	receipt_url		TEXT DEFAULT NULL,

	description		TEXT NOT NULL,
	transaction_date	DATE NOT NULL,
	created_at		TIMESTAMP DEFAULT now() NOT NULL
);

COMMIT;
