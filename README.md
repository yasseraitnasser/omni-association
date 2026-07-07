# Omni-Association
A lightweight, blazing-fast web application for managing non-profit association memberships, finances, and projects.

## Tech Stack
No heavy JavaScript frameworks, no Node.js build pipelines, and no necessary bloat. This project uses Go for runtime speed, PostgeSQL for strict relational integrity, and Server-Side Rendered HTML for a clean, reliable frontend architecture.

## Repository Blueprint
.
├── LICENSE
├── Makefile                        # High-level DevOps command abstraction
├── README.md                       # Project documentation
├── go.mod                          # Go dependency manifest
├── .env.example                    # Example of local execution configuration
├── src
│   ├── database
│   │   ├── database.go             # Go connection pool init (sql.Open + Ping)
│   │   └── migrations
│   │       └── 000_INITIAL.sql     # Base database schemas, types, and constraints
│   ├── main.go                     # Server orchestration layer & routing table
│   └── templates
│       └── layouts
│           └── layout.html         # Base HTML layout skeleton

## Local Setup & Infrastracture
1. Match Your System's Postgres Credentials
Before doing anythingm you must ensure that you local PostgreSQL instance actually has a user role that matches you configuration.
⚠️CRITICAL COTCHA: PostgreSQL will not automatically create database users based on your environment file. Whatever value you chose for `DB_USER` must exist as an authorized role in your local cluster, or the connection pool initialzation (`db.Ping()`) will fail with an authentication error.
if your chosen user doesn't exist, log into your master Postgres instance and provision them manually:
`CREATE USER "your_chosen_name" WITH PASSWORD 'your_chosen_password' CREATEDB;`
2. Configuration Environment Variables
Create a .env file at the root level of this project:
```
# Database credentials
DB_NAME=db_name
DB_PORT=5432
DB_HOST=127.0.0.1
DB_USER=user
DB_PASS=secretpass

# HTTP Server settings
PORT=8080
HOST=localhost
```
3. Execution Pipeline (Makefile Execution)
Everything else is abstracted away by the Makefile, Run `make help` (or simply `make`) to list targets, or use the code workflows below:
- Boot everything (Fresh setup): Create the DB, runs schemas, and kicks off the Go server.
```
    make all
```
- Hot-reaload Server: Starts the Go backend runtime without touching the database infrastracture.
```
    make server-run
```
- Destructive Reset: Wipes out the database entirely, builds it back from scratch according to the schema definition, and runs the backend server.
```
    make re
```
- Teardown: Drops the application database safely
```
    make server-run
```

## Relational Database Architecture
The PostgreSQL relational code handles four distinct operational domains via strict engine constranits and perfomance types:
- Strict Types: Explicit status and authorization control mapping (`member_role`, `project_status`, `project_roles`, `funding_source`, `transaction_type`, `payment_method`).
- Cascading Integrity: Key relationships (`project_members`, `project_subscribers`, `membership_payment`) leverage `ON DELETE CASCADE` actions to prevent orphan database rows during operations.
- Uniqueness Constraints: Enforces clean internal states (e.g., `unique_member_year` prevents duplicate membership fee records for a single member within a single fiscal year cycle).
