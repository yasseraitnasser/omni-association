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
	'subscriber',
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
	paid_fee	BOOLEAN DEFAULT FALSE,
	created_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS subscribers (
	id		SERIAL PRIMARY KEY,
	name		TEXT NOT NULL,
	email		VARCHAR(255) UNIQUE NOT NULL,
	paid_fee	BOOLEAN DEFAULT FALSE,
	created_at	TIMESTAMP DEFAULT now() NOT NULL
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
	role_in_project	project_roles NOT NULL DEFAULT 'contributor',
	join_at		TIMESTAMP DEFAULT now() NOT NULL,

	PRIMARY KEY (project_id, member_id)
);

CREATE TABLE IF NOT EXISTS project_subscribers (
	project_id	INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	subscriber_id	INTEGER REFERENCES subscribers(id) ON DELETE CASCADE,
	amount_paid	INTEGER NOT NULL,
	joined_at	TIMESTAMP DEFAULT now() NOT NULL,

	PRIMARY KEY (project_id, subscriber_id)
);

CREATE TABLE IF NOT EXISTS membership_payments (
	id		SERIAL PRIMARY KEY,
	member_id	INTEGER REFERENCES members(id) ON DELETE CASCADE,
	amount_paid	INTEGER NOT NULL,
	year		INTEGER NOT NULL,
	paid_at		TIMESTAMP DEFAULT now() NOT NULL,

	CONSTRAINT unique_member_year UNIQUE (member_id, year)
);

CREATE TABLE IF NOT EXISTS donations (
	id		SERIAL PRIMARY KEY,
	project_id	INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	source		funding_source NOT NULL,

	member_id	INTEGER REFERENCES members(id) ON DELETE SET NULL,
	subscriber_id	INTEGER REFERENCES subscribers(id) ON DELETE SET NULL,
	external_entity	TEXT DEFAULT NULL,

	amount_paid	INTEGER NOT NULL,
	receipt_img_url	TEXT DEFAULT NULL,
	donated_at	TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
	id			SERIAL PRIMARY KEY,
	project_id		INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	recorded_by		INTEGER REFERENCES members(id) ON DELETE SET NULL,
	type			transaction_type NOT NULL,
	amount			INTEGER NOT NULL,
	payment_method		payment_method NOT NULL,

	-- document tracking
	doc_img_url		TEXT NOT NULL,
	supplier_invoice_url	TEXT DEFAULT NULL,

	description		TEXT NOT NULL,
	transaction_date	DATE NOT NULL,
	created_at		TIMESTAMP DEFAULT now() NOT NULL
);
-- tables

COMMIT;
