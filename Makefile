## db-setup: Set up the database container
.PHONY: db-setup
db-setup:
	@cd api/sft/db && docker-compose up -d

## db-reset
.PHONY: db-reset
db-reset:
	@cd scripts && ./db-migrate-reset.sh