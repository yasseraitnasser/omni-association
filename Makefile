include .env
export
# WARN: $(MAKEFILE_LIST) <=> Makefile .env

.PHONY: help db-create db-migrate db-drop db-reset all clean re

SQL_FILE := ./src/database/migrations/000_INITIAL.sql
COLOR := \033[38;2;212;145;24m
RESET := \033[0m

help: ## Display this help and exit
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR)%-15s$(RESET) %s\n", $$1, $$2}'

all: ## Build and run the application
	$(MAKE) db-create
	$(MAKE) db-migrate
	$(MAKE) server-run

clean: ## Stop and clean the application
	$(MAKE) db-drop

re: ## Rebuild and run the application
	$(MAKE) clean
	$(MAKE) all

# backend
server-run: ## Run the backend server
	@echo 'Starting backend server...'
	@go run ./src/...
# dnekcab

# database
db-create: ## Create the database
	@echo 'Creating the database...'
	@pg_isready -h $(DB_HOST) -p $(DB_PORT) || (echo "Postgres is not running!" && exit 1)
	@PGPASSWORD=$(DB_PASS) createdb $(DB_NAME) -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER)

db-migrate: ## Migrate the database
	@echo 'Running database schema...'
	PGPASSWORD=$(DB_PASS) psql $(DB_NAME) -f $(SQL_FILE) -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER)

db-drop: ## Drop the database
	@echo 'Dropping the database...'
	@PGPASSWORD=$(DB_PASS) dropdb $(DB_NAME) -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER)

db-reset: ## Reset the database
	@echo 'Resetting the database...'
	$(MAKE) db-drop
	$(MAKE) db-create
	$(MAKE) db-migrate
# esabatad
