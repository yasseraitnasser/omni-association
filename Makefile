include .env
export
# WARN: $(MAKEFILE_LIST) <=> Makefile .env

.PHONY: help db-create db-migrate db-drop db-reset all

SQL_FILE := ./src/database/migrations/000_INITIAL.sql
COLOR := \033[38;2;212;145;24m
RESET := \033[0m

help: ## Display this help and exit
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR)%-15s$(RESET) %s\n", $$1, $$2}'

all: ## Build and start the application
	$(MAKE) db-create
	$(MAKE) db-migrate

# database begin
db-create: ## Create the database
	@echo 'Creating the database...'
	@pg_isready || (echo "Postgres is not running!" && exit 1)
	@createdb $(DB_NAME)

db-migrate: ## Migrate the database
	@echo 'Running database schema...'
	psql $(DB_NAME) -f $(SQL_FILE)

db-drop: ## Drop the database
	@echo 'Dropping the database...'
	@dropdb $(DB_NAME)

db-reset: ## Reset the database (only use if already set)
	@echo 'Resetting the database...'
	$(MAKE) db-drop
	$(MAKE) db-create
	$(MAKE) db-migrate
# database end
