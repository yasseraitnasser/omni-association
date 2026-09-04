# Omni-Association
A lightweight, blazing-fast web application for managing non-profit association memberships, finances, and projects.

## Tech Stack
No heavy JavaScript frameworks, no Node.js build pipelines, and no necessary bloat. This project uses Go for runtime speed, PostgreSQL for strict relational integrity, and Server-Side Rendered HTML for a clean, reliable frontend architecture.

## Repository Blueprint
```
.
├── LICENSE
├── Makefile                        # High-level DevOps command abstraction
├── README.md                       # Project documentation
├── go.mod                          # Go dependency manifest
├── go.sum                          # cryptographic checksum of module versions
├── .env.example                    # Example of local execution configuration
├── src
│   ├── auth
│   │   ├── admin.go                # Admin user integration
│   │   ├── jwt.go                  # JWT generation and authentication
│   │   └── login.go                # Login (validating credentials)
│   ├── database
│   │   ├── database.go             # Go connection pool init (sql.Open + Ping)
│   │   └── migrations
│   │       └── 000_INITIAL.sql     # Base database schemas, types, and constraints
│   ├── main.go                     # Server orchestration layer & routing table
│   ├── members
│   │   └── members.go              # Members management (invite, accept, decline ..)
│   ├── projects
│   │   ├── projects.go             # Project creation and management with committee assignment
│   │   ├── report.go               # Project report (summary of the project states, incomes, expenses ..)
│   │   └── transactions.go         # Project transactions tracking (income, expenses, documentations ..)
│   ├── templates
│   │   └── layouts
│   │       └── layout.html         # Base HTML layout skeleton
│   └── utils
│       ├── env.go                  # Env variables load
│       ├── password.go             # Password hashing and comparaison
│       └── validate.go             # Struct validation

```

## Local Setup & Infrastructure
1. Match Your System's Postgres Credentials
Before doing anything you must ensure that your local PostgreSQL instance actually has a user role that matches your configuration.

⚠️CRITICAL GOTCHA: PostgreSQL will not automatically create database users based on your environment file. Whatever value you chose for `DB_USER` must exist as an authorized role in your local cluster, otherwise the connection pool initialization (`db.Ping()`) will fail with an authentication error.
If your chosen user doesn't exist, log into your master Postgres instance and provision then manually:
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

SERVER_PORT=8080
SERVER_HOST=localhost

ADMIN_NAME=admin
ADMIN_EMAIL=adminemail@example.com
ADMIN_PASS=adminpass123

JWT_SECRET=4d6c669ae283e8b51ba7563009cb90562c4e37e5bef7fca915f4462006e2e72fec0c9db63482f8135f95ac2f96ce3c4ffa9205c4fff263ad94c1147443686d80
JWT_EXPIRY=24h

SECURE_TOKEN_LENGTH=32
SECURE_TOKEN_EXPIRY=24h

UPLOAD_DIR=./uploads
MAX_FILE_SIZE=2
```
3. Execution Pipeline (Makefile Execution)
Everything else is abstracted away by the Makefile, Run `make help` (or simply `make`) to list targets, or use the code workflows below:
- Boot everything (Fresh setup): Create the DB, runs schemas, and kicks off the Go server.
```
    make all
```
- Hot-reload Server: Starts the Go backend runtime without touching the database infrastructure.
```
    make server-run
```
- Destructive Reset: Wipes out the database entirely, builds it back from scratch according to the schema definition, and runs the backend server.
```
    make re
```
- Teardown: Drops the application database safely
```
    make clean
```

## Relational Database Architecture
The PostgreSQL relational code handles four distinct operational domains via strict engine constraints and performance types:
- Strict Types: Explicit status and authorization control mapping (`member_role`, `project_status`, `project_roles`, `funding_source`, `transaction_type`, `payment_method`).
- Cascading Integrity: Key relationships (`project_members`, `project_subscribers`, `membership_payment`) leverage `ON DELETE CASCADE` actions to prevent orphan database rows during operations.
- Uniqueness Constraints: Enforces clean internal states (e.g., `unique_member_year` prevents duplicate membership fee records for a single member within a single fiscal year cycle).

## Project Tasks & Roadmap
Below is the current checklist of requirements and features for this project.
- [x] The association's members are: President, Vice-President, Treasurer, Assistant Treasurer, General Secretary, Assistant General Secretary, Advisors, and members.
- [x] The association's board members are also members.
- [x] To enable the association to carry out its tasks and small projects, each member pays an annual membership fee based on the results of the general assembly, against a receipt signed and stamped by the association.
- [x] The association undertakes community projects (drilling wells, building and renovating mosques, widows' and orphans' homes, etc.).
- [x] The sources of funding for these projects are: association members (their own money or donations), government donations, and donations from other associations.
- [x] The funding source for some projects is subscriber contributions (each subscriber must pay a fee).
- [x] Before starting each project, the association's board creates a committee of members who will manage the project from start to finish.
- [x] On the platform, committee members must be able to enter transactions for each project according to their respective access rights.
- [x] The committee for each project must enter the amount of each member's donation, the payment receipt, and the supplier invoices.
- [x] At the end of each project, the committee must present its report to the association members, detailing income and expenses, along with invoices, vouchers, and copies of checks.
- [x] The association's board members must present their annual report to the subscribers once a year.
- [x] Every transaction made on the association's bank account must be documented, and a photo of the check or transfer order must be included with the reports.
- [ ] On the platform, each subscriber will be able to view their history, including their donations, subscription status, and outstanding balances owed to the association.

## Quick Progress Summary
- Completed: 12 / 13 (92%)
- Status: In Active Development
